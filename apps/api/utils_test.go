package main

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateId(t *testing.T) {
	id := GenerateId()
	assert.Len(t, id, 32)

	id2 := GenerateId()
	assert.NotEqual(t, id, id2)
}

func TestZipDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	// add random files
	err := os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("hello"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("world"), 0644)
	assert.NoError(t, err)

	zipData, err := ZipDirectory(tmpDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, zipData)

	// read zip file
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	assert.NoError(t, err)
	assert.Len(t, reader.File, 2)

	foundFiles := map[string]bool{
		"file1.txt": false,
		"file2.txt": false,
	}

	fileContents := map[string]string{
		"file1.txt": "hello",
		"file2.txt": "world",
	}

	for _, file := range reader.File {
		foundFiles[file.Name] = true

		content, err := file.Open()
		assert.NoError(t, err)
		defer content.Close()

		contentBytes, err := io.ReadAll(content)

		assert.NoError(t, err)
		assert.Equal(t, fileContents[file.Name], string(contentBytes))
	}

	assert.True(t, foundFiles["file1.txt"])
	assert.True(t, foundFiles["file2.txt"])
}
