package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UploadFile uploads a file to the specified directory
func UploadFile(file *multipart.FileHeader, directory string) (string, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(directory, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Generate a unique filename based on timestamp and file hash
	filename := generateUniqueFilename(file.Filename)

	// Create the destination file
	dst, err := os.Create(filepath.Join(directory, filename))
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return filename, nil
}

// DeleteFile deletes a file from the specified directory
func DeleteFile(filename, directory string) error {
	path := filepath.Join(directory, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	return os.Remove(path)
}

// generateUniqueFilename generates a unique filename based on timestamp and file hash
func generateUniqueFilename(originalFilename string) string {
	// Get file extension
	ext := filepath.Ext(originalFilename)
	
	// Generate hash from original filename and current timestamp
	timestamp := time.Now().UnixNano()
	hasher := md5.New()
	io.WriteString(hasher, originalFilename)
	io.WriteString(hasher, fmt.Sprintf("%d", timestamp))
	hash := hex.EncodeToString(hasher.Sum(nil))

	// Create new filename
	basename := strings.TrimSuffix(filepath.Base(originalFilename), ext)
	sanitizedBasename := sanitizeFilename(basename)
	
	return fmt.Sprintf("%s-%s%s", sanitizedBasename, hash[:8], ext)
}

// sanitizeFilename sanitizes a filename by replacing invalid characters
func sanitizeFilename(filename string) string {
	// Replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")
	
	// Remove invalid characters
	invalid := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalid {
		filename = strings.ReplaceAll(filename, char, "")
	}
	
	return filename
} 