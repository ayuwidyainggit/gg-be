package adapter

import (
	"fmt"
	"io"
	"mobile/model"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorageAdapter struct {
	BasePath   string
	BaseURL    string
	PublicPath string
}

func InitLocalStorageAdapter(basePath, baseURL, publicPath string) (*LocalStorageAdapter, error) {
	// Create base directory if not exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	// Create public directory if not exists
	if publicPath != "" {
		if err := os.MkdirAll(publicPath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create public directory: %w", err)
		}
	}

	return &LocalStorageAdapter{
		BasePath:   basePath,
		BaseURL:    baseURL,
		PublicPath: publicPath,
	}, nil
}

func (l *LocalStorageAdapter) UploadFile(req *model.Upload) (fullUrl string, err error) {
	// Generate file name
	key, err := req.GenerateFileName()
	if err != nil {
		return "", err
	}

	// Create directory structure
	fullPath := filepath.Join(l.BasePath, key)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Open source file
	src, err := req.File.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// If public path is set, also copy to public directory for serving
	if l.PublicPath != "" {
		publicFullPath := filepath.Join(l.PublicPath, key)
		publicDir := filepath.Dir(publicFullPath)
		if err := os.MkdirAll(publicDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create public directory: %w", err)
		}

		publicSrc, err := os.Open(fullPath)
		if err != nil {
			return "", fmt.Errorf("failed to reopen file for public copy: %w", err)
		}
		defer publicSrc.Close()

		publicDst, err := os.Create(publicFullPath)
		if err != nil {
			return "", fmt.Errorf("failed to create public file: %w", err)
		}
		defer publicDst.Close()

		if _, err = io.Copy(publicDst, publicSrc); err != nil {
			return "", fmt.Errorf("failed to copy to public directory: %w", err)
		}
	}

	// Generate URL
	if l.BaseURL != "" {
		// Remove leading slash from key if exists
		key = strings.TrimPrefix(key, "/")
		fullUrl = fmt.Sprintf("%s/%s", strings.TrimSuffix(l.BaseURL, "/"), key)
	} else {
		fullUrl = fmt.Sprintf("/uploads/%s", key)
	}

	return fullUrl, nil
}
