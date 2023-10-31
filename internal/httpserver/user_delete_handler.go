package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db"

	"github.com/google/uuid"
)

func (s *Server) handleDeleteUser() http.HandlerFunc {
	type request struct {
		ID string `json:"id"`
	}
	type response struct {
	}
	return func(w http.ResponseWriter, r *http.Request) {

		queries := db.New(s.DB)

		deleteUserParams := request{}

		err := json.NewDecoder(r.Body).Decode(&deleteUserParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		uuid, err := uuid.Parse(deleteUserParams.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = queries.DeleteUser(r.Context(), uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
