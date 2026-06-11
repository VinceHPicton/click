package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"

	"github.com/google/uuid"
)

func (s *Server) authAttemptCreateUserHandler() http.HandlerFunc {

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

		err := decoder.Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		confirmAuthAttemptParams := sqlc.AuthAttemptCreateUserParams{
			ID:          req.ID,
			OneTimeCode: req.OneTimeCode,
		}

		newUserID, err := s.Queries.AuthAttemptCreateUser(r.Context(), confirmAuthAttemptParams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if newUserID == uuid.Nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		accessToken, err := s.TokenManager.GenerateAccessToken(newUserID.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		refreshToken := tokens.GenerateRefreshToken()

		err = s.Queries.CreateRefreshToken(r.Context(), sqlc.CreateRefreshTokenParams{
			UserID:    newUserID,
			TokenHash: tokens.HashRefreshToken(refreshToken),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := response{
			UserID:       newUserID,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(responseBytes)
	}
}
