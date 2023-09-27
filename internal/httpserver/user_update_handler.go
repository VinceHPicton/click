package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db"
)

func (s *Server) handleUpdateUser() http.HandlerFunc {
	type request struct {
	}
	type response struct {
	}
	return func(w http.ResponseWriter, r *http.Request) {

		queries := db.New(s.DB)

		updateUserParams := db.UpdateUserParams{}

		err := json.NewDecoder(r.Body).Decode(&updateUserParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		_, err = queries.UpdateUser(r.Context(), updateUserParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
