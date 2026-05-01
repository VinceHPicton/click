package httpserver

import (
	"vincehpicton/click/internal/db"
	"encoding/json"
	"net/http"
)

type request struct {
	Mobile string `json:"mobile"`
}

func (r request) Valid() bool {
	if len(r.Mobile) == 0 {
		return false
	}
	return true
}

type response struct {
	OneTimeCode string `json:"oneTimeCode"`
}

func (s *Server) registerAttemptCreateHandler() http.HandlerFunc {
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

		responseBytes, err := json.Marshal(registerAttempt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !createRegisterAttemptParams.Valid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Write(responseBytes)
		// w.WriteHeader(http.StatusOK)
	}
}
