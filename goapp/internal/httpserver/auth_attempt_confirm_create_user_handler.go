package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"

	"github.com/google/uuid"
)

func (s *Server) authAttemptConfirmCreateUserHandler() http.HandlerFunc {

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

		authAttempt, err := s.Queries.GetValidAuthAttempt(r.Context(), req.ID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid auth attempt ID"))
			return
		}

		// TODO: auth attempt always consumed, even if user got it wrong
		err = s.Queries.ConsumeAuthAttempt(r.Context(), authAttempt.ID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Failed to consume auth attempt"))
			return
		}

		if authAttempt.OneTimeCode != req.OneTimeCode {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid one-time code"))
			return
		}

		// If successful, query for the user by mobile
		relevantUsers, err := s.Queries.GetActiveUsersByMobile(r.Context(), authAttempt.Mobile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(relevantUsers) != 0 {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		// If no user found, create a new user - TODO: this will fail if the user is banned, so we should be checking that before here really?
		// Or we can just return something like "user is banned, or another failure occurred"
		newUser, err := s.Queries.CreateUserWithMobile(r.Context(), authAttempt.Mobile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accessToken, err := s.TokenManager.GenerateAccessToken(newUser.ID.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		refreshToken, err := tokens.GenerateRefreshToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = s.Queries.CreateRefreshToken(r.Context(), sqlc.CreateRefreshTokenParams{
			UserID:    newUser.ID,
			TokenHash: tokens.HashRefreshToken(refreshToken),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := response{
			UserID:       newUser.ID,
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
