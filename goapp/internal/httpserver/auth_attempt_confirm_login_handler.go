package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"

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
	return func(w http.ResponseWriter, r *http.Request) {

		confirmLoginParams := confirmLoginRequest{}

		err := json.NewDecoder(r.Body).Decode(&confirmLoginParams)

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Get the one-time code/login req if it exists
		authAttempt, err := s.Queries.GetAuthAttempt(r.Context(), confirmLoginParams.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if authAttempt.OneTimeCode != confirmLoginParams.OneTimeCode {
			http.Error(w, "Invalid one-time code", http.StatusUnauthorized)
			return
		}

		// Get user by mobile number
		relevantUsers, err := s.Queries.GetActiveUserByMobile(r.Context(), authAttempt.Mobile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(relevantUsers) > 1 {
			// TODO: log big error here, this should never happen
			// TODO: remove this informative error in prod
			http.Error(w, "multiple users found", http.StatusInternalServerError)
			return
		}
		if len(relevantUsers) == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		user := relevantUsers[0]

		refreshToken, err := tokens.GenerateRefreshToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hashedRefreshToken := tokens.HashRefreshToken(refreshToken)

		params := sqlc.CreateRefreshTokenParams{
			UserID:    user.ID,
			TokenHash: hashedRefreshToken,
		}

		err = s.Queries.ConsumeAuthAttemptByID(r.Context(), confirmLoginParams.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = s.Queries.CreateRefreshToken(r.Context(), params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accessToken, err := s.TokenManager.GenerateAccessToken(user.ID.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := confirmLoginResponse{
			UserID:       user.ID,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode response"))
			return
		}
	}
}
