package common

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileService struct {
	StoragePath string
	BaseURL     string
}

func NewFileService(storagePath, baseURL string) (*FileService, error) {
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &FileService{
		StoragePath: storagePath,
		BaseURL:     baseURL,
	}, nil
}

func (fs *FileService) SaveFile(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", errors.New("file is nil")
	}

	orgFileName := file.Filename
	fileExt := filepath.Ext(orgFileName)

	newFilename := uuid.New().String() + fileExt

	filePath := filepath.Join(fs.StoragePath, newFilename)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	fileURL := fs.GetFileURL(newFilename)
	return fileURL, nil
}

func (fs *FileService) GetFileURL(filename string) string {
	baseURL := fs.BaseURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return baseURL + filename
}

func (fs *FileService) ServeFile(w http.ResponseWriter, filename string) {
	filePath := filepath.Join(fs.StoragePath, filename)
	http.ServeFile(w, nil, filePath)
}

func (fs *FileService) DeleteFile(filename string) error {
	filePath := filepath.Join(fs.StoragePath, filename)
	return os.Remove(filePath)
}
