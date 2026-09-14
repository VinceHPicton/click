package service

import (
	"context"
	"fmt"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"

	"github.com/google/uuid"
)

type Session struct {
	UserID       uuid.UUID
	AccessToken  string
	RefreshToken string
}

func (s *Service) createUserWithCredentials(ctx context.Context, attempt sqlc.AppAuthAttempt) (Session, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Session{}, err
	}
	// Rollback doesnt matter once Commit has succeeded.
	defer func() { _ = tx.Rollback() }()

	queriesTx := s.queries.WithTx(tx)

	// If no user found, create a new user - TODO: this will fail if the user is banned, so we should be checking that before here really?
	// Or we can just return something like "user is banned, or another failure occurred"
	newUser, err := s.queries.CreateUserWithMobile(ctx, attempt.Mobile)
	if err != nil {
		return Session{}, err
	}
	
	session, err := issueSession(ctx, queriesTx, s.tokenMgr, newUser.ID)
	if err != nil {
		return Session{}, err
	}

	if err := tx.Commit(); err != nil {
		return Session{}, err
	}

	return session, nil
}

func issueSession(ctx context.Context, queries *sqlc.Queries, tokenMgr *tokens.Manager, userID uuid.UUID) (Session, error) {
	refreshToken, err := tokens.GenerateRefreshToken()
	if err != nil {
		return Session{}, fmt.Errorf("generate refresh token: %w", err)
	}

	accessToken, err := tokenMgr.GenerateAccessToken(userID.String())
	if err != nil {
		return Session{}, fmt.Errorf("generate access token: %w", err)
	}

	_, err = queries.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: tokens.HashRefreshToken(refreshToken),
	})
	if err != nil {
		return Session{}, fmt.Errorf("persist refresh token: %w", err)
	}

	return Session{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
