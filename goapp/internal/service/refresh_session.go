package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"
)

// RefreshSession exchanges a valid refresh token for a new token pair.
//
// Rotation happens in a single statement (RotateRefreshToken revokes the old
// row and inserts the new one), so a concurrent refresh with the same token
// finds nothing to revoke and gets ErrInvalidRefreshToken rather than a second
// valid session.
func (s *Service) RefreshSession(ctx context.Context, refreshToken string) (Session, error) {
	hash := tokens.HashRefreshToken(refreshToken)

	token, err := s.queries.GetValidTokenByHash(ctx, hash)
	if errors.Is(err, sql.ErrNoRows) {
		return Session{}, ErrInvalidRefreshToken
	}
	if err != nil {
		return Session{}, fmt.Errorf("look up refresh token: %w", err)
	}

	user, err := s.queries.GetUser(ctx, token.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		// The token outlived its user, so the token itself is what's bad.
		return Session{}, ErrInvalidRefreshToken
	}
	if err != nil {
		return Session{}, fmt.Errorf("look up user: %w", err)
	}

	if user.BannedAt.Valid {
		return Session{}, ErrUserBanned
	}

	if user.DeletedAt.Valid {
		return Session{}, ErrUserDeleted
	}

	newRefreshToken, err := tokens.GenerateRefreshToken()
	if err != nil {
		return Session{}, fmt.Errorf("generate refresh token: %w", err)
	}

	_, err = s.queries.RotateRefreshToken(ctx, sqlc.RotateRefreshTokenParams{
		UserID:       user.ID,
		OldTokenHash: hash,
		NewTokenHash: tokens.HashRefreshToken(newRefreshToken),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return Session{}, ErrInvalidRefreshToken
	}
	if err != nil {
		return Session{}, fmt.Errorf("rotate refresh token: %w", err)
	}

	accessToken, err := s.tokenMgr.GenerateAccessToken(user.ID.String())
	if err != nil {
		return Session{}, fmt.Errorf("generate access token: %w", err)
	}

	return Session{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
