// Package config holds the application configuration settings.
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the application configuration settings.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Jwt      JwtConfig
	Aws      AwsConfig
	Upload   UploadConfig
}

// ServerConfig holds the server-related configuration settings.
type ServerConfig struct {
	Port    string
	GinMode string
}

// DatabaseConfig holds the database-related configuration settings.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JwtConfig holds the JWT-related configuration settings.
type JwtConfig struct {
	SecretKey             string
	ExpiresIn             time.Duration
	RefreshTokenExpiresIn time.Duration
}

// AwsConfig holds the AWS-related configuration settings.
type AwsConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3BucketName    string
	S3Endpoint      string
}

// UploadConfig holds the file upload-related configuration settings.
type UploadConfig struct {
	Path        string
	MaxFileSize int64
}

// Load loads the configuration settings from environment variables and returns a Config struct.
func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiresIn, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "15m"))
	refreshTokenExpiresIn, _ := time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_EXPIRES_IN", "720h"))
	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)

	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5433"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "myapp"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Jwt: JwtConfig{
			SecretKey:             getEnv("JWT_SECRET_KEY", "mysecretkey"),
			ExpiresIn:             jwtExpiresIn,
			RefreshTokenExpiresIn: refreshTokenExpiresIn,
		},
		Aws: AwsConfig{
			Region:          getEnv("AWS_REGION", "us-east-1"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", "test"),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", "test"),
			S3BucketName:    getEnv("S3_BUCKET_NAME", "bitly-uploads"),
			S3Endpoint:      getEnv("S3_ENDPOINT", "http://localhost:4566"),
		},
		Upload: UploadConfig{
			Path:        getEnv("UPLOAD_PATH", "./uploads"),
			MaxFileSize: maxUploadSize,
		},
	}, nil
}

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is not present, it returns the provided default value.
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
