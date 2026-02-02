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

func TestDownSampleDataset(t *testing.T) {
	t.Run("small dataset", func(t *testing.T) {
		// Case 1: Small dataset
		smallData := make([]map[string]interface{}, 100)
		for i := 0; i < 100; i++ {
			smallData[i] = map[string]interface{}{"id": i}
		}
		result := DownSampleDataset(smallData)
		assert.Equal(t, len(smallData), len(result))
		assert.Equal(t, smallData[0]["id"], result[0]["id"])
	})

	t.Run("large dataset", func(t *testing.T) {
		// Case 2: Large dataset
		largeSize := 20000
		largeData := make([]map[string]interface{}, largeSize)
		for i := 0; i < largeSize; i++ {
			largeData[i] = map[string]interface{}{"id": i}
		}
		result := DownSampleDataset(largeData)
		assert.Equal(t, 10000, len(result))

		// Check sampling logic for 2x downsample
		// index 0 -> 0
		// index 1 -> 2
		assert.Equal(t, 0, result[0]["id"])
		assert.Equal(t, 2, result[1]["id"])
		assert.Equal(t, 19998, result[9999]["id"])
	})
}
