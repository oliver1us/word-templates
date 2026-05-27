package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateHandler(c *gin.Context) {
	// Bind to a flat map to accept parameters dynamically
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Extract the template file name
	templateVal, ok := req["template"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: template"})
		return
	}
	template, ok := templateVal.(string)
	if !ok || template == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template field must be a non-empty string"})
		return
	}

	// Validate template exists
	templatePath := filepath.Join("docs", template)
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Template %s not found", template)})
		return
	}

	// Prepare data for replacement (any key other than "template")
	replaceMap := make(map[string]string)
	for k, v := range req {
		if k == "template" {
			continue
		}
		strVal := fmt.Sprintf("%v", v)
		// Handle literal \n from JSON if it was escaped as a string
		strVal = strings.ReplaceAll(strVal, "\\n", "\n")
		replaceMap[k] = strVal
	}

	// Save to a temporary file
	tmpDocxName := fmt.Sprintf("tmp/generated_%d.docx", time.Now().UnixNano())
	
	// Perform replacement using our custom robust replacer
	err := DocxReplace(templatePath, tmpDocxName, replaceMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template: " + err.Error()})
		return
	}
	defer os.Remove(tmpDocxName) // Clean up docx

	// Convert to PDF
	pdfPath, err := ConvertToPDF(tmpDocxName, "tmp")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert to PDF: " + err.Error()})
		return
	}
	defer os.Remove(pdfPath) // Clean up pdf after sending

	// Send file
	c.FileAttachment(pdfPath, strings.Replace(template, ".docx", ".pdf", 1))
}
