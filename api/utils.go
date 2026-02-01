package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// ExtractCSVFilesFromZip extracts all CSV files from a zip archive.
// It takes the zip file contents as a byte slice and returns a map
// where keys are the filenames and values are the file contents.
// Only files with a .csv extension (case-insensitive) are extracted.
func ExtractCSVFilesFromZip(zipData []byte) (map[string][]byte, error) {
	// Create a reader from the byte slice
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, fmt.Errorf("failed to read zip archive: %w", err)
	}

	csvFiles := make(map[string][]byte)

	for _, file := range reader.File {
		// Skip directories
		if file.FileInfo().IsDir() {
			continue
		}

		// Check if the file has a .csv extension (case-insensitive)
		ext := strings.ToLower(filepath.Ext(file.Name))
		if ext != ".csv" {
			continue
		}

		// Open the file within the zip
		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s in zip: %w", file.Name, err)
		}

		// Read the file contents
		contents, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s in zip: %w", file.Name, err)
		}

		// Use the base name as the key to avoid path issues
		csvFiles[filepath.Base(file.Name)] = contents
	}

	return csvFiles, nil
}

// CSVFileToJson parses CSV data and converts it to a slice of maps.
// The first row is treated as the header row, and each subsequent row
// becomes a map with header fields as keys and cell values as values.
func CSVFileToJson(csvData []byte) ([]map[string]interface{}, error) {
	reader := csv.NewReader(bytes.NewReader(csvData))

	// Read all records from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV data: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// First row contains headers
	headers := records[0]
	if len(headers) == 0 {
		return nil, fmt.Errorf("CSV file has no headers")
	}

	// Convert remaining rows to maps
	result := make([]map[string]interface{}, 0, len(records)-1)

	for i, row := range records[1:] {
		if len(row) != len(headers) {
			return nil, fmt.Errorf("row %d has %d columns, expected %d", i+2, len(row), len(headers))
		}

		rowMap := make(map[string]interface{}, len(headers))
		for j, header := range headers {
			rowMap[header] = row[j]
		}
		result = append(result, rowMap)
	}

	return result, nil
}

// GenerateId creates a new unique identifier by generating a UUID v4
// and removing all hyphens, resulting in a 32-character hex string.
func GenerateId() string {
	id := uuid.New().String()
	return strings.ReplaceAll(id, "-", "")
}

// ZipDirectory creates a zip archive containing all files in the source directory.
// It recursively walks the directory tree and adds each file to the archive.
// Returns the zip file contents as a byte slice.
func ZipDirectory(source string) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		// zip always uses forward slashes
		relPath = filepath.ToSlash(relPath)

		f, err := w.Create(relPath)
		if err != nil {
			return err
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		_, err = io.Copy(f, srcFile)
		return err
	})

	if err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
