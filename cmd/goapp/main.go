package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"vincehpicton/click/internal/httpserver"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	host     = "dbserver"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	router := mux.NewRouter()

	server := httpserver.Server{
		DB:     conn,
		Router: router,
	}

	server.Routes()

	port := 3030
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), server.Router))

}
