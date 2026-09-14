package httpserver

import (
	"encoding/json"
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
		req := request{}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "malformed request body", http.StatusBadRequest)
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
