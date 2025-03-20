package handlers

import (
	"net/http"
	"time"
	"webapp/db"
)

func HealthCheckHandler(dbConnection db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")

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

		healthCheck := db.HealthCounter{
			Datetime: time.Now().UTC(),
		}

		if err := dbConnection.Create(&healthCheck); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
