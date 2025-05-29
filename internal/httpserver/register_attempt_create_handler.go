package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db"
)

func (s *Server) registerAttemptCreateHandler() http.HandlerFunc {
	type request struct {
		Mobile string `json:"mobile"`
	}
	type response struct {
		id string
	}
	return func(w http.ResponseWriter, r *http.Request) {

		queries := db.New(s.DB)

		createRegisterAttemptParams := request{}

		err := json.NewDecoder(r.Body).Decode(&createRegisterAttemptParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		registerAttempt, err := queries.RegisterAttemptCreate(r.Context(), createRegisterAttemptParams.Mobile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		returnJsonBytes, err := json.Marshal(registerAttempt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(returnJsonBytes)
	}
}
