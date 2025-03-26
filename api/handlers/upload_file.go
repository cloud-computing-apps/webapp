package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
	"webapp/db"
	"webapp/middleware"
)

func UploadFileHandler(dbConnection db.Database, s3Client *s3.Client, bucketName string, client *statsd.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		log.WithFields(log.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
		}).Info("Received file upload request")

		if r.Method != http.MethodPost {
			log.Warn("Method not allowed", "method", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			log.Error("Failed to retrieve file from request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileID := uuid.New()
		s3Key := fmt.Sprintf("%s/%s", fileID.String(), header.Filename)

		log.WithFields(log.Fields{
			"fileID":   fileID,
			"fileName": header.Filename,
			"s3Key":    s3Key,
		}).Info("Processing file upload")

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, file)
		if err != nil {
			log.Error("Failed to read the file", "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		err = middleware.WrapS3Call(client, "upload_file", func() error {
			_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket: &bucketName,
				Key:    &s3Key,
				Body:   bytes.NewReader(buf.Bytes()),
			})
			return err
		})
		if err != nil {
			log.Error("Failed to upload file to S3", "error", err, "s3Key", s3Key)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		fileURL := fmt.Sprintf("/%s/%s", bucketName, s3Key)

		fileRecord := db.FileTable{
			Id:         fileID,
			FileName:   header.Filename,
			Url:        fileURL,
			UploadDate: time.Now(),
		}
		err = middleware.WrapDBQuery(client, "create_file_record", func() error {
			return dbConnection.Create(&fileRecord)
		})
		if err != nil {
			log.Error("Failed to create file record in DB", "error", err, "fileRecord", fileRecord)
			err = middleware.WrapS3Call(client, "delete_file_upload", func() error {
				_, err2 := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
					Bucket: &bucketName,
					Key:    &s3Key,
				})
				return err2
			})
			if err != nil {
				log.Error("Failed to rollback S3 file deletion", "error", err, "s3Key", s3Key)
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		log.WithFields(log.Fields{
			"fileRecord": fileRecord,
		}).Info("File successfully uploaded and record created")
		response, err := json.Marshal(fileRecord)
		if err != nil {
			log.Error("Failed to marshal file record", "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	}
}
