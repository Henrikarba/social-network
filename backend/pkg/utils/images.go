package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func SaveImage(data []byte, mimeType string, location string) (string, error) {
	var fileExtension string
	switch mimeType {
	case "image/jpeg":
		fileExtension = ".jpg"
	case "image/png":
		fileExtension = ".png"
	case "image/gif":
		fileExtension = ".gif"
	default:
		return "", fmt.Errorf("unsupported MIME type: %s", mimeType)
	}

	name := uuid.New().String()
	filename := name + fileExtension
	dir := fmt.Sprintf("pkg/db/files/images/%s", location)
	path := filepath.Join(dir, filename)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return "", err
	}

	// Construct the final path without "pkg/db/files"
	finalPath := strings.TrimPrefix(path, "pkg/db/files/images/")

	return finalPath, nil
}
