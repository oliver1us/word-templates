package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ConvertToPDF uses LibreOffice in headless mode to convert a docx to pdf.
func ConvertToPDF(inputFile string, outDir string) (string, error) {
	// The command: libreoffice --headless --convert-to pdf <file> --outdir <dir>
	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "pdf", inputFile, "--outdir", outDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("libreoffice error: %s, output: %s", err.Error(), string(output))
	}

	// Construct the expected PDF path
	baseName := filepath.Base(inputFile)
	pdfName := strings.Replace(baseName, ".docx", ".pdf", 1)
	pdfPath := filepath.Join(outDir, pdfName)

	return pdfPath, nil
}
