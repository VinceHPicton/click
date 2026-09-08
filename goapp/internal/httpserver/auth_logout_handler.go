package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/tokens"
)

type logoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (s *Server) logoutHandler() http.HandlerFunc {
	type response struct {
	}
	return func(w http.ResponseWriter, r *http.Request) {

		req := logoutRequest{}

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

		_, err = s.Queries.RevokeRefreshTokenByID(r.Context(), token.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to revoke refresh token"))
			return
		}
	}
}
