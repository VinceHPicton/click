package sqlc_test

func (ts *DatabaseSuite) TestCreateAuthAttempt() {
	var err error

	attempt, err := ts.queries.AuthAttemptCreate(ts.ctx, "+447712345678")
	ts.Require().NoError(err)
	ts.Require().LessOrEqual(int32(100000), attempt.OneTimeCode)
	ts.Require().GreaterOrEqual(int32(999999), attempt.OneTimeCode)

	attempts, err := ts.queries.GetAuthAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(attempts))
}
