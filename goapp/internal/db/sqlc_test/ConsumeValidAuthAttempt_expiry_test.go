package sqlc_test

import (
	"database/sql"
	"time"

	"vincehpicton/click/internal/db/sqlc"

	"github.com/google/uuid"
)

func (ts *DatabaseSuite) consumeAndReload(id uuid.UUID) sqlc.AppAuthAttempt {
	err := ts.queries.ConsumeValidAuthAttempt(ts.ctx, id)
	ts.Require().NoError(err)

	attempt, err := ts.queries.GetAuthAttempt(ts.ctx, id)
	ts.Require().NoError(err)
	return attempt
}

func (ts *DatabaseSuite) setAuthAttemptCreatedAt(id uuid.UUID, createdAt time.Time) {
	err := ts.queries.SetAuthAttemptCreatedAt(ts.ctx, sqlc.SetAuthAttemptCreatedAtParams{
		ID:        id,
		CreatedAt: sql.NullTime{Time: createdAt, Valid: true},
	})
	ts.Require().NoError(err)
}

func (ts *DatabaseSuite) TestConsumeValidAuthAttempt_ConsumesUnexpiredAttempt() {
	attempt := ts.createAuthAttemptExpiringAt(time.Now().Add(time.Minute))

	got := ts.consumeAndReload(attempt.ID)
	ts.True(got.UsedAt.Valid)
}

func (ts *DatabaseSuite) TestConsumeValidAuthAttempt_IgnoresExpiredAttempt() {
	attempt := ts.createAuthAttemptExpiringAt(time.Now().Add(-time.Hour))

	got := ts.consumeAndReload(attempt.ID)
	ts.False(got.UsedAt.Valid)
}

func (ts *DatabaseSuite) TestConsumeValidAuthAttempt_IgnoresJustExpiredAttempt() {
	attempt := ts.createAuthAttemptExpiringAt(time.Now().Add(-time.Second))

	got := ts.consumeAndReload(attempt.ID)
	ts.False(got.UsedAt.Valid)
}

func (ts *DatabaseSuite) TestConsumeValidAuthAttempt_IgnoresExpiredAttemptWithRecentCreatedAt() {
	attempt := ts.createAuthAttemptExpiringAt(time.Now().Add(-time.Minute))
	ts.setAuthAttemptCreatedAt(attempt.ID, time.Now())

	got := ts.consumeAndReload(attempt.ID)
	ts.False(got.UsedAt.Valid)
}

func (ts *DatabaseSuite) TestConsumeValidAuthAttempt_ConsumesUnexpiredAttemptWithOldCreatedAt() {
	attempt := ts.createAuthAttemptExpiringAt(time.Now().Add(time.Minute))
	ts.setAuthAttemptCreatedAt(attempt.ID, time.Now().Add(-time.Hour))

	got := ts.consumeAndReload(attempt.ID)
	ts.True(got.UsedAt.Valid)
}

func (ts *DatabaseSuite) TestConsumeValidAuthAttempt_DoesNotRestampExpiredUsedAttempt() {
	attempt := ts.createAuthAttemptExpiringAt(time.Now().Add(-time.Hour))

	err := setAuthAttemptUsedAt(ts.ctx, ts.queries, attempt.ID, time.Now().Add(-2*time.Hour))
	ts.Require().NoError(err)

	before, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)

	after := ts.consumeAndReload(attempt.ID)
	ts.Equal(before.UsedAt.Time, after.UsedAt.Time)
}
