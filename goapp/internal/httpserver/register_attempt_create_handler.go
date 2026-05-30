package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db"

	"github.com/google/uuid"
)

type registerAttemptRequest struct {
	Mobile string `json:"mobile"`
}

func (r registerAttemptRequest) Valid() bool {
	if len(r.Mobile) == 0 {
		return false
	}
	return true
}

type registerAttemptResponse struct {
	ID          uuid.UUID `json:"id"`
	OneTimeCode string    `json:"oneTimeCode"`
}

func (s *Server) registerAttemptCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		queries := db.New(s.DB)

		createRegisterAttemptParams := registerAttemptRequest{}

		err := json.NewDecoder(r.Body).Decode(&createRegisterAttemptParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		if !createRegisterAttemptParams.Valid() {
			w.WriteHeader(http.StatusBadRequest)
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

		w.Write(responseBytes)
		// w.WriteHeader(http.StatusOK)
	}
}
