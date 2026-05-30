package httpserver

import (
	"database/sql"
	"vincehpicton/click/internal/db/sqlc"

	"github.com/gorilla/mux"
)

type Server struct {
	DB      *sql.DB
	Router  *mux.Router
	Queries *sqlc.Queries
}
