package myapp

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body := []byte("Hello World")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func dbTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	dbString := "postgres://postgres:postgres@172.18.0.2:5432/db?sslmode=disable"
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

func RESTRun(r *mux.Router) {
	r.HandleFunc("/", helloWorldHandler)
	r.HandleFunc("/dbtest", dbTestHandler)
}
