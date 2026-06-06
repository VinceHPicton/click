package httpserver

import (
	"database/sql"
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

func init() {
}
