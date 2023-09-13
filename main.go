package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	dbString    = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	validInsert = "INSERT INTO matches (id, user_1_id, user_2_id, created_at) VALUES (1, 1, 1, '2022-10-10 11:30:30')"
)

func main() {

	db, err := sql.Open("postgres", dbString)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO matches (id, user_1_id, user_2_id, created_at) VALUES (1, 1, 1, '2022-10-10 11:30:30')")

	result, err := db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(rowCount)
}

func main2() {
	router := mux.NewRouter()
	RESTRun(router)

	port := 8080
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), router))
}

func RESTRun(r *mux.Router) {
	r.HandleFunc("/", helloWorldHandler)
	r.HandleFunc("/dbtest", dbTestHandler)
}

func dbTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	dbString := "postgres://postgres:postgres@localhost:5432/db?sslmode=disable"
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO matches (1, 1, 1, '2022-10-10 11:30:30')")

	result, err := db.Exec(query)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("%v", rowCount)))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body := []byte("Hello World")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
