package handlers

import (
	"gorm.io/gorm"
	"net/http"
	"time"
	"webapp/db"
)

func HealthCheckHandler(dbConnection *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.ContentLength > 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(r.URL.Query()) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		healthCheck := db.HealthCounter{
			Datetime: time.Now(),
		}

		if err := dbConnection.Create(&healthCheck).Error; err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
