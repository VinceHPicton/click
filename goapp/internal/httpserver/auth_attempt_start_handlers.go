package httpserver

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type authAttemptStartRequest struct {
	Mobile string `json:"mobile"`
}

func (r authAttemptStartRequest) Valid() bool {
	return len(r.Mobile) > 0
}

type authAttemptStartResponse struct {
	ID          uuid.UUID `json:"id"`
	OneTimeCode string    `json:"oneTimeCode"`
}

func (s *Server) authAttemptStartLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := decodeJSON[authAttemptStartRequest](w, r)
		if !ok {
			return
		}

		attempt, err := s.Queries.AuthAttemptCreate(r.Context(), req.Mobile)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, authAttemptStartResponse{
			ID:          attempt.ID,
			OneTimeCode: strconv.Itoa(int(attempt.OneTimeCode)),
		})
	}
}

func (s *Server) authAttemptStartCreateUserHandler() http.HandlerFunc {
	// For now, literally just calling the other one, in future we will probably have these diverge in some way
	return s.authAttemptStartLoginHandler()
}
