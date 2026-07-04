package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"
)

func (s *Server) refreshHandler() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refreshToken"`
	}
	type response struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
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

		hash := tokens.HashRefreshToken(req.RefreshToken)
		token, err := s.Queries.GetValidTokenByHash(r.Context(), hash)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid refresh token"))
			return
		}

		user, err := s.Queries.GetUser(r.Context(), token.UserID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("User not found"))
			return
		}

		newRefreshToken := tokens.GenerateRefreshToken()
		p := sqlc.RotateRefreshTokenParams{
			UserID:       user.ID,
			OldTokenHash: hash,
			NewTokenHash: tokens.HashRefreshToken(newRefreshToken),
		}

		_, err = s.Queries.RotateRefreshToken(r.Context(), p)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Failed to rotate refresh token"))
			return
		}

		accessToken, err := s.TokenManager.GenerateAccessToken(user.ID.String())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Failed to generate access token"))
			return
		}

		resp := response{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
