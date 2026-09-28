package sqlc_test

import (
	"database/sql"
	"errors"
	"vincehpicton/click/internal/db/factory"

	"github.com/google/uuid"
)

const (
	phoneNumber = "+447712345678"
)

func (ts *DatabaseSuite) TestCreateAuthAttempt() {
	var err error

	attempt, err := factory.FakeAuthAttemptWithServiceExpiry(ts.ctx, ts.queries, phoneNumber)
	ts.Require().NoError(err)
	ts.Require().LessOrEqual(int32(100000), attempt.OneTimeCode)
	ts.Require().GreaterOrEqual(int32(999999), attempt.OneTimeCode)

	attempts, err := ts.queries.GetAuthAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(attempts))
}

func (ts *DatabaseSuite) TestGetAuthAttempt() {

	var err error

	attempt, err := factory.FakeAuthAttemptWithServiceExpiry(ts.ctx, ts.queries, phoneNumber)
	ts.Require().NoError(err)

	getAttempt, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(attempt.ID, getAttempt.ID)
}

func (ts *DatabaseSuite) TestGetValidAuthAttempt() {
	var err error

	attempt, err := factory.FakeAuthAttemptWithServiceExpiry(ts.ctx, ts.queries, phoneNumber)
	ts.Require().NoError(err)

	getAttempt, err := ts.queries.GetValidAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(attempt.ID, getAttempt.ID)
}

func (ts *DatabaseSuite) TestGetValidAuthAttempt_DoesntExist() {
	_, err := ts.queries.GetValidAuthAttempt(ts.ctx, uuid.New())
	ts.Require().Error(err)

	ts.True(errors.Is(err, sql.ErrNoRows))
}

func (ts *DatabaseSuite) TestConsumeValidAuthAttempt() {
	var err error

	attempt, err := factory.FakeAuthAttemptWithServiceExpiry(ts.ctx, ts.queries, phoneNumber)
	ts.Require().NoError(err)

	err = ts.queries.ConsumeValidAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)

	consumedAttempt, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(consumedAttempt.UsedAt.Valid, true)
}

func (ts *DatabaseSuite) TestSQLCExecDoesNotErrorIfNoRowsAffected() {
	var err error

	err = ts.queries.ConsumeValidAuthAttempt(ts.ctx, uuid.New())
	ts.Require().NoError(err)
}
