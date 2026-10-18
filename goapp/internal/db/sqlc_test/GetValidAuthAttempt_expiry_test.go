package sqlc_test

import (
	"database/sql"
	"errors"
	"time"

	"vincehpicton/click/internal/db/sqlc"
)

func (ts *DatabaseSuite) createAuthAttemptExpiringAt(expiresAt time.Time) sqlc.AppAuthAttempt {
	attempt, err := ts.queries.AuthAttemptCreate(ts.ctx, sqlc.AuthAttemptCreateParams{
		Mobile:    phoneNumber,
		ExpiresAt: expiresAt,
	})
	ts.Require().NoError(err)
	return attempt
}

func (ts *DatabaseSuite) TestGetValidAuthAttempt_ReturnsUnexpiredAttempt() {
	attempt := ts.createAuthAttemptExpiringAt(time.Now().Add(time.Minute))

	got, err := ts.queries.GetValidAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.Equal(attempt.ID, got.ID)
}

func (ts *DatabaseSuite) TestGetValidAuthAttempt_RejectsExpiredAttempt() {
	attempt := ts.createAuthAttemptExpiringAt(time.Now().Add(-time.Hour))

	_, err := ts.queries.GetValidAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().Error(err)
	ts.True(errors.Is(err, sql.ErrNoRows))
}

func (ts *DatabaseSuite) TestGetValidAuthAttempt_RejectsJustExpiredAttempt() {
	attempt := ts.createAuthAttemptExpiringAt(time.Now().Add(-time.Second))

	_, err := ts.queries.GetValidAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().Error(err)
	ts.True(errors.Is(err, sql.ErrNoRows))
}
