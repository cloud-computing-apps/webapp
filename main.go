package main

import (
	"net/http"
	"webapp/api"
	"webapp/db"
	"webapp/setup"
)

func main() {
	dbInstance := setup.DBConn()
	db.InitDB(dbInstance)

	s3Client := setup.S3Conn()

	routes := api.RegisterRoutes(dbInstance, s3Client, setup.S3Bucket)

	err := http.ListenAndServe(":8080", routes)
	if err != nil {
		panic(err)
	}
}
