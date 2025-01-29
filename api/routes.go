package api

import (
	"gorm.io/gorm"
	"net/http"
	"webapp/api/handlers"
)

func RegisterRoutes(dbConnection *gorm.DB) http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/healthz", handlers.HealthCheckHandler(dbConnection))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, pattern := r.Handler(req)
		if pattern == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		r.ServeHTTP(w, req)
	})
}
