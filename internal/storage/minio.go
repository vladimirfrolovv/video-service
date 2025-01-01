package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"

	"github.com/vladimirfrolovv/video-service/internal/config"
)

func NewMinioClient(minioCfg config.MinioConfig) (*minio.Client, error) {
	minioClient, err := minio.New(minioCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioCfg.AccessKey, minioCfg.SecretKey, ""),
		Secure: minioCfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}

func EnsureBucket(client *minio.Client, minioCfg config.MinioConfig) error {
	ctx := context.Background()
	exists, errBucketExists := client.BucketExists(ctx, minioCfg.BucketName)
	if errBucketExists != nil {
		return fmt.Errorf("ошибка при проверке бакета %s: %w", minioCfg.BucketName, errBucketExists)
	}
	if !exists {
		if err := client.MakeBucket(ctx, minioCfg.BucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("не удалось создать бакет %s: %w", minioCfg.BucketName, err)
		}
	}
	log.Printf("Бакет %s уже существует.\n", minioCfg.BucketName)
	return nil
}

func ListObjects(client *minio.Client, bucketName string) ([]string, error) {
	ctx := context.Background()

	var files []string
	objectCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	for obj := range objectCh {
		if obj.Err != nil {
			return nil, obj.Err
		}
		files = append(files, obj.Key)
	}

	return files, nil
}
