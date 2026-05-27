package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Ensure docs and tmp directories exist
	os.MkdirAll("docs", os.ModePerm)
	os.MkdirAll("tmp", os.ModePerm)

	r.POST("/generate", GenerateHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Starting server on port %s...", port)
	r.Run(":" + port)
}
