package config

import (
	"os"
)

type Config struct {
	AppPort string
	Minio   MinioConfig
}

type MinioConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	UseSSL     bool
	BucketName string
}

func LoadConfig() *Config {
	return &Config{
		AppPort: getEnv("APP_PORT", ":8080"),
		Minio: MinioConfig{
			Endpoint:   getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey:  getEnv("MINIO_ACCESS_KEY", "admin"),
			SecretKey:  getEnv("MINIO_SECRET_KEY", "admin"),
			UseSSL:     getEnv("MINIO_USE_SSL", "false") == "true",
			BucketName: getEnv("MINIO_BUCKET_NAME", "videos"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
