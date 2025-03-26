package handlers

import (
	"github.com/DataDog/datadog-go/statsd"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	"webapp/db"
	"webapp/middleware"
)

func HealthCheckHandler(dbConnection db.Database, client *statsd.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
		}).Info("Received health check request")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		if r.Method != http.MethodGet {
			log.Warn("Method not allowed", "method", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.ContentLength > 0 {
			log.Warn("Bad request: content length should be zero", "contentLength", r.ContentLength)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(r.URL.Query()) > 0 {
			log.Warn("Bad request: query parameters not allowed", "query", r.URL.Query())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		healthCheck := db.HealthCounter{
			Datetime: time.Now().UTC(),
		}

		err := middleware.WrapDBQuery(client, "create_health_counter", func() error {
			return dbConnection.Create(&healthCheck)
		})
		if err != nil {
			log.Error("Failed to create health counter", "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		log.WithFields(log.Fields{
			"datetime": healthCheck.Datetime,
		}).Debug("Health check counter updated")

		log.WithFields(log.Fields{
			"datetime": healthCheck.Datetime,
		}).Info("Health check recorded successfully")

		w.WriteHeader(http.StatusOK)
	}
}
