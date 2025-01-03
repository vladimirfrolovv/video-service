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
		log.Println(".env file not read", err)
	}
	cfg := config.LoadConfig()

	minioClient, err := storage.NewMinioClient(cfg.Minio)
	if err != nil {
		log.Fatalf("Not initialize minio client: %v\n", err)
	}
	if err := storage.EnsureBucket(minioClient, cfg.Minio); err != nil {
		log.Fatalf("Dont check exists minio bucket: %v\n", err)
	}

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(simpleCORS)
	router.HandleFunc("/upload", handlers.UploadHandler(minioClient, cfg.Minio)).Methods("POST")
	router.HandleFunc("/video/{filename}", handlers.GetVideoHandler(minioClient, cfg.Minio)).Methods("GET")
	router.HandleFunc("/list", handlers.ListHandler(minioClient, cfg.Minio)).Methods("GET")
	fmt.Println("Сервис запущен на порту", cfg.AppPort)
	if err := http.ListenAndServe(cfg.AppPort, router); err != nil {
		log.Fatal(err)
	}
}
func simpleCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		//Only options
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
