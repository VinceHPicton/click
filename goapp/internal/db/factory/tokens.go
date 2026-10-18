package factory

import (
	"context"
	"time"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"

	"github.com/google/uuid"
)

type refreshTokenParams struct {
	expiresAt *time.Time
}

type RefreshTokenOption func(*refreshTokenParams)

func FakeRefreshToken(
	ctx context.Context,
	q *sqlc.Queries,
	userID uuid.UUID,
	options ...RefreshTokenOption,
) (unhashedToken string, dbToken sqlc.AppRefreshToken, err error) {
	params := refreshTokenParams{}

	for _, option := range options {
		option(&params)
	}

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
	if err != nil {
		return refreshToken, token, err
	}

	if params.expiresAt == nil {
		return refreshToken, token, nil
	}

	err = q.SetRefreshTokenExpiryByID(ctx, sqlc.SetRefreshTokenExpiryByIDParams{
		ID:        token.ID,
		ExpiresAt: *params.expiresAt,
	})
	if err != nil {
		return refreshToken, token, err
	}

	token, err = q.GetTokenByHash(ctx, refreshTokenHash)

	return refreshToken, token, err
}

func WithExpiry(expiresAt time.Time) RefreshTokenOption {
	return func(p *refreshTokenParams) {
		p.expiresAt = &expiresAt
	}
}
