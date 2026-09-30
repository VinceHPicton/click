package httpserver

import (
	"net/http"

	"github.com/google/uuid"
)

type confirmLoginRequest struct {
	ID          uuid.UUID `json:"id"`
	OneTimeCode int32     `json:"oneTimeCode"`
}

type confirmLoginResponse struct {
	UserID       uuid.UUID `json:"userId"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}

func (s *Server) authAttemptConfirmLoginHandler() http.HandlerFunc {
	service := s.service()

	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := decodeJSON[confirmLoginRequest](w, r)
		if !ok {
			return
		}

		session, err := service.ConfirmLogin(r.Context(), req.ID, req.OneTimeCode)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, confirmLoginResponse{
			UserID:       session.UserID,
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
		})
	}
}
