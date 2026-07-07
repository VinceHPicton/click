package httpserver

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) authAttemptConfirmLoginHandler() http.HandlerFunc {
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

	}
}
