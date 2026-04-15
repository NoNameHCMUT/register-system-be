package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func Save(file multipart.File, header *multipart.FileHeader, subdir, uploadDir string) (string, error) {
	if err := os.MkdirAll(filepath.Join(uploadDir, subdir), 0o755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return "", fmt.Errorf("unsupported file type: %s", contentType)
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".jpg"
	}

	filename := uuid.New().String() + ext
	dst := filepath.Join(uploadDir, subdir, filename)

	dstFile, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return "/" + subdir + "/" + filename, nil
}
