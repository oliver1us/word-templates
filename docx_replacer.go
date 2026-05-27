package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// DocxReplace replaces placeholders in a docx file.
// It supports {{key}} with optional spaces (e.g. {{ key }} or {{key}}).
func DocxReplace(inputPath, outputPath string, replaceMap map[string]string) error {
	// 1. Open the zip reader
	reader, err := zip.OpenReader(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input zip: %w", err)
	}
	defer reader.Close()

	// 2. Create the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// 3. Create the zip writer
	writer := zip.NewWriter(outFile)
	defer writer.Close()

	// 4. Iterate through files in the zip
	for _, file := range reader.File {
		// Create a writer for this file in the new zip
		w, err := writer.Create(file.Name)
		if err != nil {
			return fmt.Errorf("failed to create file %s in output: %w", file.Name, err)
		}

		// Open the source file
		r, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s in input: %w", file.Name, err)
		}

		// Check if it's an XML file inside the word folder (where document text is stored)
		if strings.HasSuffix(file.Name, ".xml") && strings.Contains(file.Name, "word/") {
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			r.Close()
			if err != nil {
				return fmt.Errorf("failed to read xml content: %w", err)
			}

			content := buf.String()

			// Replace each key in the map
			for k, v := range replaceMap {
				// Escape XML special characters in the replacement value
				escapedVal := xmlEscape(v)

				// Handle newlines: replace \n with </w:t><w:br/><w:t> in Word XML
				escapedVal = strings.ReplaceAll(escapedVal, "\n", "</w:t><w:br/><w:t>")

				// Regex to match {{ k }} or {{k}} (case insensitive)
				escapedKey := regexp.QuoteMeta(k)
				re, err := regexp.Compile(`(?i)\{\{\s*` + escapedKey + `\s*\}\}`)
				if err == nil {
					// Use ReplaceAllLiteralString to avoid dollar sign ($) replacement bugs in Go Regex
					content = re.ReplaceAllLiteralString(content, escapedVal)
				}
			}

			// Write the modified content
			_, err = io.WriteString(w, content)
			if err != nil {
				return fmt.Errorf("failed to write modified xml content: %w", err)
			}
		} else {
			// Copy non-XML files directly
			_, err = io.Copy(w, r)
			r.Close()
			if err != nil {
				return fmt.Errorf("failed to copy file %s: %w", file.Name, err)
			}
		}
	}

	return nil
}

// xmlEscape escapes special XML characters.
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
