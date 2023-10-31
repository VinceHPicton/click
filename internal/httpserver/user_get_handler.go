package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db"

	"github.com/google/uuid"
)

func (s *Server) handleGetUser() http.HandlerFunc {
	type request struct {
		ID string `json:"id"`
	}
	type response struct {
	}
	return func(w http.ResponseWriter, r *http.Request) {

		queries := db.New(s.DB)

		getUserParams := request{}

		err := json.NewDecoder(r.Body).Decode(&getUserParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		uuid, err := uuid.Parse(getUserParams.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := queries.GetUser(r.Context(), uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userJsonBytes, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(userJsonBytes)

		w.WriteHeader(http.StatusOK)
	}
}
