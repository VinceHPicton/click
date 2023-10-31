package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db"
)

func (s *Server) userUpdateHandler() http.HandlerFunc {
	type request struct {
	}
	type response struct {
	}
	return func(w http.ResponseWriter, r *http.Request) {

		queries := db.New(s.DB)

		userUpdateParams := db.UserUpdateParams{}

		err := json.NewDecoder(r.Body).Decode(&userUpdateParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		_, err = queries.UserUpdate(r.Context(), userUpdateParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
