package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db"
)

func (s *Server) handleCreateUser() http.HandlerFunc {
	type request struct {
	}
	type response struct {
	}
	return func(w http.ResponseWriter, r *http.Request) {

		queries := db.New(s.DB)

		createUserParams := db.CreateUserParams{}

		err := json.NewDecoder(r.Body).Decode(&createUserParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		user, err := queries.CreateUser(r.Context(), createUserParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte(user.ID.String()))
	}
}
