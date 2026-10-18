package service

import (
	"database/sql"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"
)

type Service struct {
	queries  *sqlc.Queries
	tokenMgr *tokens.Manager
	db       *sql.DB
}

func New(db *sql.DB, queries *sqlc.Queries, tokenManager *tokens.Manager) *Service {
	return &Service{
		db:      db,
		queries: queries,
		tokenMgr:  tokenManager,
	}
}
