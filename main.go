package main

import (
	"github.com/codescratchers/golang-webserver/api"
	"github.com/codescratchers/golang-webserver/database"
	"log"
	"net/http"
)

func main() {
	db, err := database.ConnectToMySQL("hushtabs", "hushtabs", "127.0.0.1:3306", "hushtabs_db")

	if err != nil {
		log.Fatal(err)
		return
	}

	// initialize apiServer
	server := api.NewApiServer(":8080", database.Storage{DB: db}, http.NewServeMux())

	// start server
	api.Serve(server)
}
