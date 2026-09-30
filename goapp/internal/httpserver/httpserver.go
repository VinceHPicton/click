package httpserver

import (
	"database/sql"
	"net/http"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/service"
	"vincehpicton/click/internal/tokens"

	"github.com/gorilla/mux"
)

type Server struct {
	DB           *sql.DB
	Router       *mux.Router
	Queries      *sqlc.Queries
	TokenManager *tokens.Manager

	// Service is optional. When nil, service builds one from the fields above,
	// so a plain Server literal remains a valid way to construct the server.
	Service *service.Service
}

// service returns the service, building a default one if none was
// injected. Callers should invoke this when constructing a handler rather than
// per request.
func (s *Server) service() *service.Service {
	if s.Service != nil {
		return s.Service
	}

	return service.New(s.DB, s.Queries, s.TokenManager)
}

func (s *Server) authMiddleware() func(http.Handler) http.Handler {
	return AuthMiddleware(s.TokenManager)
}

func (s *Server) CORSMiddleware() func(http.Handler) http.Handler {
	return CORSMiddleware()
}
