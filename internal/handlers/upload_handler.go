package handlers

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/minio/minio-go/v7"

	"github.com/vladimirfrolovv/video-service/internal/config"
)

func UploadHandler(minioClient *minio.Client, minioCfg config.MinioConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ограничение на размер загружаемого файла — 100 МБ
		if err := r.ParseMultipartForm(100 << 20); err != nil {
			http.Error(w, "Ошибка при чтении multipart формы: "+err.Error(), http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Не удалось получить файл из запроса: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		uploadInfo, err := uploadToMinIO(r.Context(), minioClient, minioCfg.BucketName, file, handler)
		if err != nil {
			log.Printf("Ошибка загрузки файла: %v\n", err)
			http.Error(w, "Ошибка при загрузке файла в MinIO: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Файл успешно загружен. ETag: %s\n", uploadInfo.ETag)
	}
}

func uploadToMinIO(ctx context.Context, client *minio.Client, bucketName string, file multipart.File, handler *multipart.FileHeader) (minio.UploadInfo, error) {
	objectName := handler.Filename

	if _, err := file.Seek(0, 0); err != nil {
		return minio.UploadInfo{}, err
	}
	return client.PutObject(
		ctx,
		bucketName,
		objectName,
		file,
		handler.Size,
		minio.PutObjectOptions{
			// invalid part size
			//PartSize:    500,
			NumThreads:  4,
			ContentType: handler.Header.Get("Content-Type"),
		},
	)
}
