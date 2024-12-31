package handlers

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"

	"github.com/vladimirfrolovv/video-service/internal/config"
)

func GetVideoHandler(minioClient *minio.Client, minioCfg config.MinioConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename, ok := vars["filename"]
		if !ok || filename == "" {
			http.Error(w, "filename не указан", http.StatusBadRequest)
			return
		}

		object, err := minioClient.GetObject(context.Background(), minioCfg.BucketName, filename, minio.GetObjectOptions{})
		if err != nil {
			log.Printf("Не удалось получить объект %s из MinIO: %v\n", filename, err)
			http.Error(w, "Ошибка при получении видео", http.StatusInternalServerError)
			return
		}
		defer object.Close()

		_, err = object.Stat()
		if err != nil {
			log.Printf("Не удалось выполнить Stat() для объекта %s: %v\n", filename, err)
			http.Error(w, "Видео не найдено или ошибка чтения", http.StatusNotFound)
			return
		}

		_, err = io.Copy(w, object)
		if err != nil {
			log.Printf("Ошибка при копировании данных объекта %s: %v\n", filename, err)
		}
	}
}
