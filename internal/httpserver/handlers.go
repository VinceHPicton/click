package httpserver

import (
	"fmt"
	"net/http"
)

func (s *Server) handleInsertTest() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		query := "INSERT INTO matches (id, user_1_id, user_2_id, created_at) VALUES ($1, $2, $3, $4)"

		statement, err := s.DB.Prepare(query)
		if err != nil {
			fmt.Println(err.Error())
			w.Write([]byte(err.Error()))
			return
		}
		defer statement.Close()

		result, err := statement.Exec(1, 1, 1, "2022-10-10 11:30:30")
		if err != nil {
			fmt.Println(err.Error())
			w.Write([]byte(err.Error()))
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			fmt.Println(err.Error())
			w.Write([]byte(err.Error()))
			return
		}

		rowsAffBytes := []byte(fmt.Sprintf("%v", rowsAffected))

		w.Write(rowsAffBytes)

	}
}

func (s *Server) handleDBping() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		err := s.DB.Ping()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte("DB pinged"))
	}
}
