package factory

import (
	"context"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"

	"github.com/google/uuid"
)

func FakeRefreshToken(ctx context.Context, q *sqlc.Queries, userID uuid.UUID) (sqlc.AppRefreshToken, error) {
	refreshToken := tokens.GenerateRefreshToken()
	refreshTokenHash := tokens.HashRefreshToken(refreshToken)

	p := sqlc.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: refreshTokenHash,
	}

	return q.CreateRefreshToken(ctx, p)
}
