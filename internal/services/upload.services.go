package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/zekielmp/Bitly/internal/config"
	"github.com/zekielmp/Bitly/internal/interfaces"
)

type UploadService struct {
	provider interfaces.UploadProvider
}

func NewUploadService(provider interfaces.UploadProvider) *UploadService {
	return &UploadService{
		provider: provider,
	}
}

func (s *UploadService) UploadProductImage(productID uint, file *multipart.FileHeader) (string, error) {

	ext := strings.ToLower(filepath.Ext(file.Filename))
	newFile := uuid.New().String() + ext
	if os.Getenv("UPLOAD_PROVIDER") == "cld" {
		// Remove the file extension from the filename
		path := fmt.Sprintf("products/%d/%s", productID, newFile)
		return s.provider.UploadFile(file, path)
	}

	if !isvalidImageExt(ext) {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	path := fmt.Sprintf("products/%d/%s%s", productID, newFile, ext)
	return s.provider.UploadFile(file, path)
}

func isvalidImageExt(ext string) bool {
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

func (s *UploadService) cloudinary(cfg *config.Config) bool {
	return cfg.Upload.UploadProvider == "cld"
}
