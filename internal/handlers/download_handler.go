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
			http.Error(w, "filename not exists", http.StatusBadRequest)
			return
		}

		object, err := minioClient.GetObject(context.Background(), minioCfg.BucketName, filename, minio.GetObjectOptions{})
		if err != nil {
			log.Printf("Dont get object %s from minio: %v\n", filename, err)
			http.Error(w, "Dont get video", http.StatusInternalServerError)
			return
		}
		defer object.Close()

		_, err = object.Stat()
		if err != nil {
			log.Printf("Not execute Stat() for object %s: %v\n", filename, err)
			http.Error(w, "Video not exists", http.StatusNotFound)
			return
		}

		_, err = io.Copy(w, object)
		if err != nil {
			log.Printf("Dont must copy object %s: %v\n", filename, err)
		}
	}
}
