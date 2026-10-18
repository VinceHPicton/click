package service

import (
	"context"
	"database/sql"
	"errors"
	"vincehpicton/click/internal/db/sqlc"

	"github.com/google/uuid"
)

func (s *Service) ConfirmCreateUser(ctx context.Context, attemptID uuid.UUID, oneTimeCode int32) (Session, error) {
	authAttempt, err := getAndConsumeAttempt(ctx, s.queries, attemptID, oneTimeCode)
	if err != nil {
		return Session{}, err
	}

	// TODO: if user banned - this will error with the first err - maybe we should do another check first? OR use another query/edit this
	exists, err := checkActiveUserExists(ctx, s.queries, authAttempt.Mobile)
	if err != nil {
		return Session{}, err
	}
	if exists {
		return Session{}, ErrUserExists
	}

	session, err := s.createUserWithCredentials(ctx, authAttempt)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func checkActiveUserExists(ctx context.Context, queries *sqlc.Queries, mobile string) (bool, error) {
	_, err := queries.GetActiveUserByMobile(ctx, mobile)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, err
}
