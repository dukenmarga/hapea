package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Document struct {
	Latex string `uri:"latex" binding:"required"`
}

func main() {
	r := gin.Default()
	r.GET("/api/v1/simple/:latex", func(c *gin.Context) {
		var doc Document
		if err := c.ShouldBindUri(&doc); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		id := uuid.New()
		c.JSON(http.StatusOK, gin.H{
			"id":    id.String(),
			"latex": doc.Latex,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
