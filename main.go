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

	routes := api.RegisterRoutes(dbInstance)

	err := http.ListenAndServe(":8080", routes)
	if err != nil {
		panic(err)
	}
}
