package httpserver

import (
	"database/sql"
	"net/http"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"

	"github.com/gorilla/mux"
)

type Server struct {
	DB           *sql.DB
	Router       *mux.Router
	Queries      *sqlc.Queries
	TokenManager *tokens.Manager
}

func (s *Server) authMiddleware() func(http.Handler) http.Handler {
	return AuthMiddleware(s.TokenManager)
}

func (s *Server) CORSMiddleware() func(http.Handler) http.Handler {
	return CORSMiddleware()
}
