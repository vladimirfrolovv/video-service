package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"github.com/vladimirfrolovv/video-service/internal/config"
	"github.com/vladimirfrolovv/video-service/internal/handlers"
	"github.com/vladimirfrolovv/video-service/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Внимание: .env файл не найден или не может быть прочитан:", err)
	}
	cfg := config.LoadConfig()

	minioClient, err := storage.NewMinioClient(cfg.Minio)
	if err != nil {
		log.Fatalf("Не удалось инициализировать MinIO клиент: %v\n", err)
	}
	if err := storage.EnsureBucket(minioClient, cfg.Minio); err != nil {
		log.Fatalf("Не удалось убедиться в существовании бакета: %v\n", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/upload", handlers.UploadHandler(minioClient, cfg.Minio)).Methods("POST")
	fmt.Println("Сервис запущен на порту", cfg.AppPort)
	if err := http.ListenAndServe(cfg.AppPort, router); err != nil {
		log.Fatal(err)
	}
}
