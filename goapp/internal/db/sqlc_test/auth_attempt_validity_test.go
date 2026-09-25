package sqlc_test

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/service"

	"github.com/google/uuid"
)

func (ts *DatabaseSuite) TestGetValidAuthAttempt_AcceptsAttemptInsideWindow() {
	attempt, err := ts.queries.AuthAttemptCreate(ts.ctx, phoneNumber)
	ts.Require().NoError(err)

	fiveSecondsInsideWindowTime := time.Now().Add(-service.AuthAttemptValidityWindowMinutes*time.Minute).Add(time.Second*5)

	err = setAuthAttemptCreatedAt(ts.ctx, ts.queries, attempt.ID, fiveSecondsInsideWindowTime)
	ts.Require().NoError(err)

	validAttempt, err := ts.queries.GetValidAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.Equal(attempt.ID, validAttempt.ID)
}

func (ts *DatabaseSuite) TestGetValidAuthAttempt_RejectsAttemptOutsideWindow() {
	attempt, err := ts.queries.AuthAttemptCreate(ts.ctx, phoneNumber)
	ts.Require().NoError(err)

	fiveSecondsBehindWindowTime := time.Now().Add(-service.AuthAttemptValidityWindowMinutes*time.Minute).Add(-time.Second*5)

	err = setAuthAttemptCreatedAt(ts.ctx, ts.queries, attempt.ID, fiveSecondsBehindWindowTime)
	ts.Require().NoError(err)

	_, err = ts.queries.GetValidAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().Error(err)
	ts.True(errors.Is(err, sql.ErrNoRows))

	stillThere, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.False(stillThere.UsedAt.Valid)
}

func (ts *DatabaseSuite) TestGetValidAuthAttempt_RejectsAlreadyUsedAttempt() {
	attempt, err := ts.queries.AuthAttemptCreate(ts.ctx, phoneNumber)
	ts.Require().NoError(err)

	err = ts.queries.ConsumeAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)

	_, err = ts.queries.GetValidAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().Error(err)
	ts.True(errors.Is(err, sql.ErrNoRows))
}

// ConsumeAuthAttempt is the query the service calls. Its used_at IS NULL guard
// is what makes consumption idempotent, so a replay cannot refresh the stamp
// and buy another two minutes of validity.
func (ts *DatabaseSuite) TestConsumeAuthAttempt_RefusesAlreadyUsedAttempt() {
	attempt, err := ts.queries.AuthAttemptCreate(ts.ctx, phoneNumber)
	ts.Require().NoError(err)

	err = setAuthAttemptUsedAt(ts.ctx, ts.queries, attempt.ID, time.Now().Add(-time.Hour))
	ts.Require().NoError(err)

	before, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.Require().True(before.UsedAt.Valid)

	err = ts.queries.ConsumeAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)

	after, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.Equal(before.UsedAt.Time, after.UsedAt.Time, "an already used attempt must not be re-stamped")
}

func (ts *DatabaseSuite) TestConsumeAuthAttempt_HasNoFreshnessGuard() {
	attempt, err := ts.queries.AuthAttemptCreate(ts.ctx, phoneNumber)
	ts.Require().NoError(err)

	oneHourAgo := time.Now().Add(-time.Hour)

	err = setAuthAttemptCreatedAt(ts.ctx, ts.queries, attempt.ID, oneHourAgo)
	ts.Require().NoError(err)

	err = ts.queries.ConsumeAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)

	consumed, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.NotEqual(oneHourAgo, consumed.UsedAt.Time, "an attempt outside the window will still be stamped as used by ConsumeAuthAttempt")
}

func (ts *DatabaseSuite) TestConsumeAuthAttemptByID_RefusesAlreadyUsedAttempt() {
	attempt, err := ts.queries.AuthAttemptCreate(ts.ctx, phoneNumber)
	ts.Require().NoError(err)

	oneHourAgo := time.Now().Add(-time.Hour)

	err = setAuthAttemptUsedAt(ts.ctx, ts.queries, attempt.ID, oneHourAgo)
	ts.Require().NoError(err)

	before, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)

	err = ts.queries.ConsumeAuthAttemptByID(ts.ctx, attempt.ID)
	ts.Require().NoError(err)

	after, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.Equal(before.UsedAt.Time, after.UsedAt.Time)
}

func (ts *DatabaseSuite) TestConsumeAuthAttemptByID_RefusesAttemptOutsideWindow() {
	attempt, err := ts.queries.AuthAttemptCreate(ts.ctx, phoneNumber)
	ts.Require().NoError(err)

	fiveSecondsBehindWindowTime := time.Now().Add(-service.AuthAttemptValidityWindowMinutes*time.Minute).Add(-time.Second*5)

	err = setAuthAttemptCreatedAt(ts.ctx, ts.queries, attempt.ID, fiveSecondsBehindWindowTime)
	ts.Require().NoError(err)

	err = ts.queries.ConsumeAuthAttemptByID(ts.ctx, attempt.ID)
	ts.Require().NoError(err)

	untouched, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.False(untouched.UsedAt.Valid)
}

// Offsets are taken from the created_at Postgres assigned rather than from
// time.Now, so the window assertions do not depend on the Go process and the
// database container agreeing on the current time.
func setAuthAttemptCreatedAt(ctx context.Context, q *sqlc.Queries, id uuid.UUID, createdAt time.Time) error {
	return q.SetAuthAttemptCreatedAt(ctx, sqlc.SetAuthAttemptCreatedAtParams{
		ID:        id,
		CreatedAt: sql.NullTime{Time: createdAt, Valid: true},
	})
}

func setAuthAttemptUsedAt(ctx context.Context, q *sqlc.Queries, id uuid.UUID, usedAt time.Time) error {
	return q.SetAuthAttemptUsedAt(ctx, sqlc.SetAuthAttemptUsedAtParams{
		ID:     id,
		UsedAt: sql.NullTime{Time: usedAt, Valid: true},
	})
}
