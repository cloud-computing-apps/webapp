package api

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"webapp/api/handlers"
	"webapp/db"
)

func RegisterRoutes(dbConnection *gorm.DB, s3Client *s3.Client, bucketName string) http.Handler {
	r := http.NewServeMux()
	dbInstance := &db.GormDatabase{DB: dbConnection}

	r.HandleFunc("/healthz", handlers.HealthCheckHandler(dbInstance))

	r.HandleFunc("/v1/file", handlers.UploadFileHandler(dbInstance, s3Client, bucketName))

	r.HandleFunc("/v1/file/", func(w http.ResponseWriter, r *http.Request) {
		fileID := strings.TrimPrefix(r.URL.Path, "/v1/file/")
		if fileID == "" {
			http.Error(w, "Missing file id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			handlers.GetFileHandler(dbConnection, fileID)(w, r)
		case http.MethodDelete:
			// Update DeleteFileHandler to accept fileID as a parameter.
			handlers.DeleteFileHandler(dbConnection, s3Client, bucketName, fileID)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, pattern := r.Handler(req)
		if pattern == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		r.ServeHTTP(w, req)
	})
}
