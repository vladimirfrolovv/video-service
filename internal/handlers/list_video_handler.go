package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"

	"github.com/vladimirfrolovv/video-service/internal/config"
	"github.com/vladimirfrolovv/video-service/internal/storage"
)

func ListHandler(client *minio.Client, minioCfg config.MinioConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := storage.ListObjects(client, minioCfg.BucketName)
		if err != nil {
			log.Println("Error when get list video:", err)
			http.Error(w, "Error when get list video", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}
