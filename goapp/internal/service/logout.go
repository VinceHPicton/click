package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"vincehpicton/click/internal/tokens"
)

// Logout revokes the refresh token so it can no longer be rotated.
//
// The already-issued access token stays valid until it expires; revoking it
// would need a deny list, which the app does not have yet.
//
// Banned and soft-deleted users are allowed to log out. Refusing would leave
// their refresh token live, which is the opposite of what is wanted.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	hash := tokens.HashRefreshToken(refreshToken)

	token, err := s.queries.GetValidTokenByHash(ctx, hash)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrInvalidRefreshToken
	}
	if err != nil {
		return fmt.Errorf("look up refresh token: %w", err)
	}

	if _, err := s.queries.RevokeRefreshTokenByID(ctx, token.ID); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}
