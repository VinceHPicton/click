package main

import (
	"fmt"
	"log"
	"net/http"

	"vincehpicton/click/internal/db"
	"vincehpicton/click/internal/httpserver"

	"github.com/gorilla/mux"
)

func main() {

	db, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	server := httpserver.Server{
		DB:     db,
		Router: router,
	}

	server.Routes()

	port := 3030
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), server.Router))

}
