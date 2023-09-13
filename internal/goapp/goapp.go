package goapp

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

const (
	query = "INSERT INTO matches (id, user_1_id, user_2_id, created_at) VALUES ($1, $2, $3, $4)"

	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
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

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	////////////////////////////////////////
	// SECURE: using prepared statements
	statement, err := db.Prepare(query)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	// .Exec() would be user-given data here
	result, err := statement.Exec(1, 1, 1, "2022-10-10 11:30:30")
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	////////////////////////////////////////
	// NOT SECURE: using raw db.Exec
	result, err = db.Exec(query, 1, 1, 1, "2022-10-10 11:30:30")
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
