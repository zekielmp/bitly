package providers

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/zekielmp/Bitly/internal/config"
)

type CloudinaryProvider struct {
	Client    *cloudinary.Cloudinary
	CloudName string
	ApiKey    string
	ApiSecret string
}

func NewCloudinaryProvider(cfg *config.Config) *CloudinaryProvider {

	cld, err := cloudinary.NewFromParams(cfg.Cloudinary.CloudName, cfg.Cloudinary.ApiKey, cfg.Cloudinary.ApiSecret)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Cloudinary client: %v", err.Error()))
	}
	return &CloudinaryProvider{
		Client:    cld,
		CloudName: cfg.Cloudinary.CloudName,
		ApiKey:    cfg.Cloudinary.ApiKey,
		ApiSecret: cfg.Cloudinary.ApiSecret,
	}

}

func (p *CloudinaryProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() {
		if err := src.Close(); err != nil {
			return
		}
	}()

	/* Upload the file to Cloudinary */
	upload, err := p.Client.Upload.Upload(context.TODO(), src, uploader.UploadParams{
		PublicID: path,
	})

	if err != nil {
		return "", err
	}
	log.Printf("cloud=%q key=%q secretLen=%d", p.CloudName, p.ApiKey, len(p.ApiSecret))

	if upload.Error.Message != "" {
		return "", fmt.Errorf("Cloudinary upload error: %s", upload.Error.Message)
	}
	return upload.SecureURL, nil
}

func (p *CloudinaryProvider) DeleteFile(path string) error {
	/* Implement the logic to delete the file from Cloudinary here.
	You can use the Cloudinary Go SDK or make HTTP requests to the Cloudinary API.*/
	publicID := path /* Assuming the path is the public ID of the file in Cloudinary.*/
	_, err := p.Client.Upload.Destroy(context.TODO(), uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return err
	}

	// For demonstration purposes, let's assume the deletion is successful and return nil.
	return nil
}
