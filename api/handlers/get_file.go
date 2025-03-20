package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"webapp/db"
)

func GetFileHandler(dbConnection *gorm.DB, fileID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		dbInstance := &db.GormDatabase{DB: dbConnection}

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id, err := uuid.Parse(fileID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var fileRecord db.FileTable
		if err := dbInstance.FindByID(id, &fileRecord); err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		jsonResponse, _ := json.Marshal(fileRecord)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}
