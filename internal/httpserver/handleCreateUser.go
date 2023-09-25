package httpserver

import (
	"net/http"
)

func (s *Server) handleCreateUser() http.HandlerFunc {
	type request struct {
	}
	type response struct {
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// s.DB.CreateUser()
	}
}
