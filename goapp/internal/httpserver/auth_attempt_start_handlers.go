package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type authAttemptStartRequest struct {
	Mobile string `json:"mobile"`
}

func (r authAttemptStartRequest) Valid() bool {
	if len(r.Mobile) == 0 {
		return false
	}
	return true
}

type authAttemptStartResponse struct {
	ID          uuid.UUID `json:"id"`
	OneTimeCode string    `json:"oneTimeCode"`
}

func (s *Server) authAttemptStartLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		createAuthAttemptParams := authAttemptStartRequest{}

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

		resp := authAttemptStartResponse{
			ID:          authAttempt.ID,
			OneTimeCode: strconv.Itoa(int(authAttempt.OneTimeCode)),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode response"))
			return
		}
	}
}

func (s *Server) authAttemptStartCreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For now, literally just calling the other one, in future we will probably have these diverge in some way
		s.authAttemptStartLoginHandler()(w, r)
	}
}
