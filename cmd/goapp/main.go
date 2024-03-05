package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"vincehpicton/click/internal/httpserver"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

var (
	// host     = "dbserver"
	// port     = 5432
	// user     = "postgres"
	// password = "postgres"
	// dbname   = "postgres"

	dbhost     = os.Getenv("DB_ADDR")
	dbport     = os.Getenv("DB_PORT")
	dbuser     = os.Getenv("DB_USER")
	dbpassword = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpassword, dbname)

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

	appPort := os.Getenv("GOAPP_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", appPort), server.Router))

}
