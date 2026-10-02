package providers

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	appconfig "github.com/zekielmp/Bitly/internal/config"
)

type S3Provider struct {
	client     *s3.Client
	uploader   *manager.Uploader
	bucketName string
	endpoint   string
}

func NewS3Provider(cfg *appconfig.Config) *S3Provider {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Aws.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.Aws.AccessKeyID,
			cfg.Aws.SecretAccessKey,
			"")),
	)
	if err != nil {
		panic("unable to load AWS SDK config, " + err.Error())
	}

	//configure for localstack
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Aws.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Aws.S3Endpoint)
			o.UsePathStyle = true
		}
	})

	return &S3Provider{
		client:     client,
		uploader:   manager.NewUploader(client),
		bucketName: cfg.Aws.S3BucketName,
		endpoint:   cfg.Aws.S3Endpoint,
	}
}

func (p *S3Provider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}

	defer src.Close()

	dst, err := p.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(path),
		Body:   src,
	})

	if err != nil {
		return "", err
	}
	return *dst.Key, nil
}

func (p *S3Provider) DeleteFile(path string) error {
	_, err := p.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(strings.TrimPrefix(path, "/")),
	})

	return err
}
