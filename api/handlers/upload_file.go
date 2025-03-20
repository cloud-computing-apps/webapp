package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"time"
	"webapp/db"
)

func UploadFileHandler(dbConnection db.Database, s3Client *s3.Client, bucketName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileID := uuid.New()
		s3Key := fmt.Sprintf("%s/%s", fileID.String(), header.Filename)

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, file)
		if err != nil {
			log.Println("failed to read the file", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &bucketName,
			Key:    &s3Key,
			Body:   bytes.NewReader(buf.Bytes()),
		})
		if err != nil {
			log.Println("failed to add file to s3", err)
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

		if err := dbConnection.Create(&fileRecord); err != nil {
			log.Println("failed to create record", err)
			_, err2 := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
				Bucket: &bucketName,
				Key:    &s3Key,
			})
			if err2 != nil {
				log.Println("failed to delete file to s3", err2)
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		response, err := json.Marshal(fileRecord)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	}
}
