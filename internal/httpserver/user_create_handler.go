package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db"
)

func (s *Server) userCreateHandler() http.HandlerFunc {
	type request struct {
	}
	type response struct {
		id string
	}
	return func(w http.ResponseWriter, r *http.Request) {

		queries := db.New(s.DB)

		userCreateParams := db.UserCreateParams{}

		err := json.NewDecoder(r.Body).Decode(&userCreateParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		user, err := queries.UserCreate(r.Context(), userCreateParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		userJsonBytes, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(userJsonBytes)
	}
}
