package providers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalUploadProvider struct {
	basePath string
}

func NewLocalUploadProvider(basePath string) *LocalUploadProvider {
	return &LocalUploadProvider{
		basePath: basePath,
	}
}

func (p *LocalUploadProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {

	filePath := filepath.Join(p.basePath, path)

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return " ", err
	}

	/*Open source of string*/
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() {
		if err := src.Close(); err != nil {
			return
		}
	}()

	/*create destination*/
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := dst.Close(); err != nil {
			return
		}
	}()

	/*read from source to destination*/
	if _, err := dst.ReadFrom(src); err != nil {
		return "", nil
	}
	return fmt.Sprintf("/Uploads/%s", path), nil

}

func (p *LocalUploadProvider) DeleteFile(path string) error {
	fullPath := filepath.Join(p.basePath, path)
	return os.Remove(fullPath)
}
