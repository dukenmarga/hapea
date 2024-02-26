package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Document struct {
	Latex  string `uri:"latex" binding:"required"`
	Format string `uri:"format" default:"raw"` // raw, base64
}

func main() {
	r := gin.Default()

	// Simple latex conversion from URI (GET)
	r.GET("/api/v1/simple/:latex", func(c *gin.Context) {
		var err error
		id := uuid.New()

		// Map URL Query param into struct
		var doc Document
		if err := c.ShouldBindUri(&doc); err != nil {
			c.JSON(400, gin.H{
				"name": "map query",
				"id":   id.String(),
				"msg":  err,
			})
			return
		}

		tmplFile := "simple-article.tmpl.tex"
		filename := id.String()
		outputPath := "./files/"
		texFilename := filename + ".tex"
		pdfFilename := filename + ".pdf"

		// Prepare tex file handler
		f, err := openTexFile(outputPath, texFilename)
		if err != nil {
			log.Printf("Error: %s", err)
			errorResponse(c, id, "open tex", err)
			return
		}
		defer f.Close()

		// Insert latex string based on template
		tmpl, err := insertTexToTemplate(tmplFile)
		if err != nil {
			log.Printf("Error: %s", err)
			errorResponse(c, id, "parse tmpl", err)
			return
		}

		// Write the latex string into a file
		err = writeLatex(tmpl, f, doc)
		if err != nil {
			log.Printf("Error: %s", err)
			errorResponse(c, id, "write tex", err)
			return
		}

		// Convert the file tex into PDF
		err = convertTexToPDF(outputPath, texFilename)
		if err != nil {
			log.Printf("Error: %s", err)
			errorResponse(c, id, "convert tex", err)
			return
		}

		// Attach the PDF (downloaded in browser)
		attachPDF(c, outputPath, pdfFilename)

		// Delete the file
		err = deleteTmpFiles(filename, outputPath)
		if err != nil {
			log.Printf("Error: %s", err)
			return
		}
	})

	// Handle latex conversion from POST
	r.POST("/api/v1/convert", func(c *gin.Context) {
		var err error
		id := uuid.New()
		doc := Document{
			Latex:  c.PostForm("latex"),
			Format: c.PostForm("format"),
		}

		// If the format is base64, decode it
		if doc.Format == "base64" {
			var decodedByte, _ = base64.StdEncoding.DecodeString(doc.Latex)
			doc.Latex = string(decodedByte)
		}

		tmplFile := "simple-article.tmpl.tex"
		filename := id.String()
		outputPath := "./files/"
		texFilename := filename + ".tex"
		pdfFilename := filename + ".pdf"

		// Prepare tex file handler
		f, err := openTexFile(outputPath, texFilename)
		if err != nil {
			log.Printf("Error: %s", err)
			errorResponse(c, id, "open tex", err)
			return
		}
		defer f.Close()

		// Insert latex string based on template
		tmpl, err := insertTexToTemplate(tmplFile)
		if err != nil {
			log.Printf("Error: %s", err)
			errorResponse(c, id, "parse tmpl", err)
			return
		}

		// Write the latex string into a file
		err = writeLatex(tmpl, f, doc)
		if err != nil {
			log.Printf("Error: %s", err)
			errorResponse(c, id, "write tex", err)
			return
		}

		// Convert the file tex into PDF
		err = convertTexToPDF(outputPath, texFilename)
		if err != nil {
			log.Printf("Error: %s", err)
			errorResponse(c, id, "convert tex", err)
			return
		}

		// Attach the PDF (downloaded in browser)
		attachPDF(c, outputPath, pdfFilename)

		// Delete the file
		err = deleteTmpFiles(filename, outputPath)
		if err != nil {
			log.Printf("Error: %s", err)
			return
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func openTexFile(outputPath string, texFilename string) (*os.File, error) {
	f, err := os.Create(outputPath + texFilename)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func insertTexToTemplate(tmplFile string) (*template.Template, error) {
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func writeLatex(tmpl *template.Template, f *os.File, doc Document) error {
	err := tmpl.Execute(f, doc)
	if err != nil {
		return err
	}
	return nil
}

func convertTexToPDF(outputPath string, texFilename string) error {
	// pdflatex <filename.tex>
	proc := exec.Command(
		"lualatex",
		"-halt-on-error",
		"-output-directory",
		outputPath,
		outputPath+texFilename,
	)
	proc.Env = os.Environ()
	proc.Env = append(proc.Env, "buf_size=1000000")

	out := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	proc.Stdout = out
	proc.Stderr = stderr
	err := proc.Run()
	if err != nil {
		return errors.New(err.Error() + ": " + stderr.String())
	}
	return nil
}

func attachPDF(c *gin.Context, outputPath string, pdfFilename string) {
	c.FileAttachment(fmt.Sprintf("%s/%s", outputPath, pdfFilename), pdfFilename)
	c.Writer.Header().Set("Content-type", "application/octet-stream")
}

func deleteTmpFiles(filename string, outputPath string) error {
	deleteExts := []string{
		".pdf",
		".tex",
		".log",
		".aux",
	}
	for _, ext := range deleteExts {
		err := os.Remove(outputPath + filename + ext)
		if err != nil {
			return err
		}
	}
	return nil
}

func errorResponse(c *gin.Context, id uuid.UUID, name string, err error) {
	c.JSON(500, gin.H{
		"name": name,
		"id":   id.String(),
		"msg":  err.Error(),
	})
}
