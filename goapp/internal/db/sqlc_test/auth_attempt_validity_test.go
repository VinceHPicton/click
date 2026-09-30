package sqlc_test

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"vincehpicton/click/internal/db/factory"
	"vincehpicton/click/internal/db/sqlc"

	"github.com/google/uuid"
)

func (ts *DatabaseSuite) TestGetValidAuthAttempt_RejectsAlreadyUsedAttempt() {
	attempt, err := factory.FakeAuthAttemptWithServiceExpiry(ts.ctx, ts.queries, phoneNumber)
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
	attempt, err := factory.FakeAuthAttemptWithServiceExpiry(ts.ctx, ts.queries, phoneNumber)
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

func (ts *DatabaseSuite) TestConsumeValidAuthAttempt_RefusesAlreadyUsedAttempt() {
	attempt, err := factory.FakeAuthAttemptWithServiceExpiry(ts.ctx, ts.queries, phoneNumber)
	ts.Require().NoError(err)

	oneHourAgo := time.Now().Add(-time.Hour)

	err = setAuthAttemptUsedAt(ts.ctx, ts.queries, attempt.ID, oneHourAgo)
	ts.Require().NoError(err)

	before, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)

	err = ts.queries.ConsumeValidAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)

	after, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.Equal(before.UsedAt.Time, after.UsedAt.Time)
}

func setAuthAttemptUsedAt(ctx context.Context, q *sqlc.Queries, id uuid.UUID, usedAt time.Time) error {
	return q.SetAuthAttemptUsedAt(ctx, sqlc.SetAuthAttemptUsedAtParams{
		ID:     id,
		UsedAt: sql.NullTime{Time: usedAt, Valid: true},
	})
}
