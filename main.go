package main

import (
	"github.com/DataDog/datadog-go/statsd"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"webapp/api"
	"webapp/db"
	"webapp/middleware"
	"webapp/setup"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	stats, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		log.Fatal("Could not connect to StatsD: %v", err)
	}
	defer stats.Close()

	dbInstance := setup.DBConn()
	db.InitDB(dbInstance)

	s3Client := setup.S3Conn()

	routes := api.RegisterRoutes(dbInstance, s3Client, setup.S3Bucket, stats)

	wrappedRoutes := middleware.MetricMiddleware(stats, routes)

	err2 := http.ListenAndServe(":8080", wrappedRoutes)
	if err2 != nil {
		log.Fatal("Server failed: ", err2)
	}
}
