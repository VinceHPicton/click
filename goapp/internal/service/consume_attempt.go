package service

import (
	"context"
	"database/sql"
	"errors"
	"vincehpicton/click/internal/db/sqlc"

	"github.com/google/uuid"
)

func getAndConsumeAttempt(ctx context.Context, queries *sqlc.Queries, attemptID uuid.UUID, oneTimeCode int32) (sqlc.AppAuthAttempt, error) {
	attempt, err := queries.GetValidAuthAttempt(ctx, attemptID)
	if errors.Is(err, sql.ErrNoRows) {
		return sqlc.AppAuthAttempt{}, ErrInvalidAttempt
	}
	if err != nil {
		return sqlc.AppAuthAttempt{}, err
	}

	// TODO: auth attempt always consumed, even if user got it wrong, fixable with a "retries" column or something
	if err := queries.ConsumeAuthAttempt(ctx, attempt.ID); err != nil {
		return sqlc.AppAuthAttempt{}, err
	}

	if attempt.OneTimeCode != oneTimeCode {
		return sqlc.AppAuthAttempt{}, ErrInvalidCode
	}

	return attempt, nil
}
