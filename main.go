package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Document struct {
	Latex string `uri:"latex" binding:"required"`
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

		// Prepare tex file handler
		// f, err := os.OpenFile(texPath+texFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		f, err := os.Create(outputPath + texFilename)
		if err != nil {
			c.JSON(500, gin.H{
				"name": "tex handler",
				"id":   id.String(),
				"msg":  err.Error(),
			})
			return
		}

		// Insert latex string based on template
		tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
		if err != nil {
			c.JSON(500, gin.H{
				"name": "read tmpl",
				"id":   id.String(),
				"msg":  err.Error(),
			})
			return
		}

		// Write the latex string into a file
		err = tmpl.Execute(f, doc)
		if err != nil {
			c.JSON(500, gin.H{
				"name": "write tex",
				"id":   id.String(),
				"msg":  err.Error(),
			})
			return
		}

		// Convert the file tex into PDF
		// pdflatex <filename.tex>
		pdfFilename := filename + ".pdf"
		proc := exec.Command(
			"pdflatex",
			"-halt-on-error",
			"-output-directory",
			outputPath,
			outputPath+texFilename,
		)
		out := bytes.NewBuffer([]byte{})
		proc.Stdout = out
		err = proc.Run()
		if err != nil {
			c.JSON(500, gin.H{
				"name": "convert",
				"id":   id.String(),
				"msg":  out.String() + err.Error(),
			})
			return
		}

		// Attach the PDF (downloaded in browser)
		c.FileAttachment(fmt.Sprintf("%s/%s", outputPath, pdfFilename), pdfFilename)
		c.Writer.Header().Set("Content-type", "application/octet-stream")

		// Delete the file
		deleteExts := []string{
			".pdf",
			".tex",
			".log",
			".aux",
		}
		for _, ext := range deleteExts {
			err = os.Remove(outputPath + filename + ext)
			if err != nil {
				c.JSON(500, gin.H{
					"name": "delete " + ext,
					"id":   id.String(),
					"msg":  err.Error(),
				})
				return
			}
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
