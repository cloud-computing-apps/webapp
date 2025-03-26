package handlers

import (
	"encoding/json"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"webapp/db"
	"webapp/middleware"
)

func GetFileHandler(dbConnection *gorm.DB, fileID string, client *statsd.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		log.WithFields(log.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
		}).Info("Received request to get file record")

		if r.Method != http.MethodGet {
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
		err = middleware.WrapDBQuery(client, "find_file_record", func() error {
			return dbInstance.FindByID(id, &fileRecord)
		})
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Warn("File record not found", "fileID", fileID)
				w.WriteHeader(http.StatusNotFound)
			} else {
				log.Error("Failed to find file record", "fileID", fileID, "error", err)
				w.WriteHeader(http.StatusServiceUnavailable)
			}
			return
		}
		log.WithFields(log.Fields{
			"fileID":     fileID,
			"fileRecord": fileRecord,
		}).Info("Successfully retrieved file record")

		jsonResponse, err := json.Marshal(fileRecord)
		if err != nil {
			log.Error("Failed to marshal file record", "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}
