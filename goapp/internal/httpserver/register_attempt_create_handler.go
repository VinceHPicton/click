package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type authAttemptRequest struct {
	Mobile string `json:"mobile"`
}

func (r authAttemptRequest) Valid() bool {
	if len(r.Mobile) == 0 {
		return false
	}
	return true
}

type authAttemptResponse struct {
	ID          uuid.UUID `json:"id"`
	OneTimeCode string    `json:"oneTimeCode"`
}

func (s *Server) authAttemptCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		createAuthAttemptParams := authAttemptRequest{}

		err := json.NewDecoder(r.Body).Decode(&createAuthAttemptParams)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		if !createAuthAttemptParams.Valid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		authAttempt, err := s.Queries.AuthAttemptCreate(r.Context(), createAuthAttemptParams.Mobile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseBytes, err := json.Marshal(authAttempt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(responseBytes)
		// w.WriteHeader(http.StatusOK)
	}
}
