package httpserver

import (
	"net/http"
)

type logoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (s *Server) logoutHandler() http.HandlerFunc {
	service := s.service()

	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := decodeJSON[logoutRequest](w, r)
		if !ok {
			return
		}

		if err := service.Logout(r.Context(), req.RefreshToken); err != nil {
			writeError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
