package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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
	newFile := uuid.New().String()
	_, ok := isvalidImageExt(ext)
	if !ok {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}
	path := fmt.Sprintf("products/%d/%s", productID, newFile)

	return s.provider.UploadFile(file, path)
}

func isvalidImageExt(ext string) (string, bool) {
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, validExt := range validExts {
		if ext == validExt {
			return ext, true
		}
	}
	return "", false
}

func (s *UploadService) cloudinary(cfg *config.Config) bool {
	return cfg.Upload.UploadProvider == "cld"
}

func removeBG(file io.Reader, filename string) (io.ByteReader, error) {
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	go func() {
		part, err := mw.CreateFormFile("image_file", filename)
		if err == nil {
			_, err = io.Copy(part, file)
		}
		if err == nil {
			err = mw.WriteField("size", "auto")
		}
		mw.Close()
		pw.CloseWithError(err)
	}()

	req, err := http.NewRequest("POST", "https://api.remove.bg/v1.0/removebg", pr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", os.Getenv("REMOVE_BG_KEY"))
	req.Header.Set("Content-Type", mw.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("remove.bg %d: %s", res.StatusCode, body)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
