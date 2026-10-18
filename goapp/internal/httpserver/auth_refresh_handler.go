package httpserver

import (
	"net/http"
)

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type refreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (s *Server) refreshHandler() http.HandlerFunc {
	service := s.service()

	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := decodeJSON[refreshRequest](w, r)
		if !ok {
			return
		}

		session, err := service.RefreshSession(r.Context(), req.RefreshToken)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, refreshResponse{
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
		})
	}
}
