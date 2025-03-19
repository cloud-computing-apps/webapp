package handlers

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"webapp/db"
)

func DeleteFileHandler(dbConnection *gorm.DB, s3Client *s3.Client, bucketName string, fileID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id, err := uuid.Parse(fileID)
		if err != nil {
			http.Error(w, `{"error": "Invalid file ID format"}`, http.StatusBadRequest)
			return
		}

		dbInstance := &db.GormDatabase{DB: dbConnection}
		var fileRecord db.FileTable
		if err := dbInstance.FindByID(id, &fileRecord); err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s3Key := strings.TrimPrefix(fileRecord.Url, fmt.Sprintf("/%s/", bucketName))

		if _, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: &bucketName,
			Key:    &s3Key,
		}); err != nil {
			http.Error(w, `{"error": "Failed to delete file from S3"}`, http.StatusInternalServerError)
			return
		}

		if err := dbConnection.Delete(&fileRecord).Error; err != nil {
			http.Error(w, `{"error": "Failed to delete file from DB"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
