package sqlc_test

const (
	phoneNumber = "+447712345678"
)

func (ts *DatabaseSuite) TestCreateAuthAttempt() {
	var err error

	attempt, err := ts.queries.AuthAttemptCreate(ts.ctx, phoneNumber)
	ts.Require().NoError(err)
	ts.Require().LessOrEqual(int32(100000), attempt.OneTimeCode)
	ts.Require().GreaterOrEqual(int32(999999), attempt.OneTimeCode)

	attempts, err := ts.queries.GetAuthAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(attempts))
}

func (ts *DatabaseSuite) TestGetAuthAttempt() {
	var err error

	attempt, err := ts.queries.AuthAttemptCreate(ts.ctx, phoneNumber)
	ts.Require().NoError(err)

	getAttempt, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(attempt.ID, getAttempt.ID)
}

func (ts *DatabaseSuite) TestConsumeAuthAttemptByID() {
	var err error

	attempt, err := ts.queries.AuthAttemptCreate(ts.ctx, phoneNumber)
	ts.Require().NoError(err)

	err = ts.queries.ConsumeAuthAttemptByID(ts.ctx, attempt.ID)
	ts.Require().NoError(err)

	consumedAttempt, err := ts.queries.GetAuthAttempt(ts.ctx, attempt.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(consumedAttempt.UsedAt.Valid, true)
}