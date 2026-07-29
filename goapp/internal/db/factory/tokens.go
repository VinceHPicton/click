package factory

import (
	"context"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"

	"github.com/google/uuid"
)

func FakeRefreshToken(ctx context.Context, q *sqlc.Queries, userID uuid.UUID) (unhashedToken string, dbToken sqlc.AppRefreshToken, err error) {
	refreshToken, err := tokens.GenerateRefreshToken()
	if err != nil {
		return "", sqlc.AppRefreshToken{}, err
	}
	refreshTokenHash := tokens.HashRefreshToken(refreshToken)

	p := sqlc.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: refreshTokenHash,
	}

	token, err := q.CreateRefreshToken(ctx, p)

	return refreshToken, token, err
}
