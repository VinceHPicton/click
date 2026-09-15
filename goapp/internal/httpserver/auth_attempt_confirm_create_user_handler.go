package httpserver

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) authAttemptConfirmCreateUserHandler() http.HandlerFunc {
	service := s.service()

	type request struct {
		ID          uuid.UUID `json:"id"`
		OneTimeCode int32     `json:"oneTimeCode"`
	}

	type response struct {
		UserID       uuid.UUID `json:"userId"`
		AccessToken  string    `json:"accessToken"`
		RefreshToken string    `json:"refreshToken"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := decodeJSON[request](w, r)
		if !ok {
			return
		}

		session, err := service.ConfirmCreateUser(r.Context(), req.ID, req.OneTimeCode)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, response{
			UserID:       session.UserID,
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
		})
	}
}
