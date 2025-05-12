package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UploadDirectory is the directory where uploads are stored
const UploadDirectory = "uploads"

// GenerateFileName generates a unique filename for an uploaded file
func GenerateFileName(originalName string) (string, error) {
	// Get file extension
	ext := filepath.Ext(originalName)

	// Generate random bytes for filename
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return "", fmt.Errorf("failed to generate random filename: %v", err)
	}

	// Create timestamp prefix
	timestamp := time.Now().Format("20060102-150405")

	// Create a unique filename with timestamp and random hex
	filename := fmt.Sprintf("%s-%s%s", timestamp, hex.EncodeToString(randBytes), ext)

	return filename, nil
}

// SaveUploadedFile saves an uploaded file to the filesystem
func SaveUploadedFile(file *multipart.FileHeader, directory string) (string, error) {
	// Ensure the upload directory exists
	uploadPath := filepath.Join(directory, UploadDirectory)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Generate a unique filename
	filename, err := GenerateFileName(file.Filename)
	if err != nil {
		return "", err
	}

	// Construct the full path
	filePath := filepath.Join(uploadPath, filename)

	// Save the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Return the relative path from the upload directory
	return filepath.Join(UploadDirectory, filename), nil
}

// DeleteFile deletes a file from the filesystem
func DeleteFile(filePath string, directory string) error {
	// Ensure the path is not empty and is within the uploads directory
	if filePath == "" {
		return nil // Nothing to delete
	}

	// Check if the path starts with the upload directory
	if !strings.HasPrefix(filePath, UploadDirectory) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// Construct the full path
	fullPath := filepath.Join(directory, filePath)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	// Delete the file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}
