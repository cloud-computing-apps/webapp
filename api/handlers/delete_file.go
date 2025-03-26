package handlers

import (
	"context"
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"webapp/db"
	"webapp/middleware"
)

func DeleteFileHandler(dbConnection *gorm.DB, s3Client *s3.Client, bucketName string, fileID string, client *statsd.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		log.WithFields(log.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
		}).Info("Received request to delete file")

		if r.Method != http.MethodDelete {
			log.Warn("Method not allowed", "method", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id, err := uuid.Parse(fileID)
		if err != nil {
			log.Error("Failed to parse fileID", "fileID", fileID, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dbInstance := &db.GormDatabase{DB: dbConnection}
		var fileRecord db.FileTable
		if err := dbInstance.FindByID(id, &fileRecord); err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Warn("File record not found", "fileID", fileID)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			log.Error("Error retrieving file record", "fileID", fileID, "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		s3Key := strings.TrimPrefix(fileRecord.Url, fmt.Sprintf("/%s/", bucketName))
		log.WithFields(log.Fields{
			"fileID": fileID,
			"s3Key":  s3Key,
		}).Info("File record found; proceeding with deletion")

		backup := fileRecord

		err = middleware.WrapDBQuery(client, "delete_file_record", func() error {
			return dbInstance.Delete(&fileRecord)
		})
		if err != nil {
			log.Error("Failed to delete file record from DB", "fileID", fileID, "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		log.WithField("fileID", fileID).Info("File record deleted from DB")

		err = middleware.WrapS3Call(client, "delete_s3_object", func() error {
			_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
				Bucket: &bucketName,
				Key:    &s3Key,
			})
			return err
		})
		if err != nil {
			log.Error("Failed to delete object from S3", "s3Key", s3Key, "error", err)
			rollbackErr := middleware.WrapDBQuery(client, "rollback_delete_file_record", func() error {
				return dbInstance.Create(&backup)
			})
			if rollbackErr != nil {
				log.Error("Failed to rollback DB deletion", "fileID", fileID, "error", rollbackErr)
			} else {
				log.Info("Rolled back DB deletion successfully", "fileID", fileID)
			}
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		log.WithField("fileID", fileID).Info("File successfully deleted from S3 and DB")
		w.WriteHeader(http.StatusNoContent)
	}
}
