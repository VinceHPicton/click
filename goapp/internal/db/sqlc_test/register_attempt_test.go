package sqlc_test

func (ts *DatabaseSuite) TestRegisterAttemptCreate() {
	registerAttempt, err := ts.queries.RegisterAttemptCreate(ts.ctx, "+447840195455")
	ts.NoError(err)
	ts.NotEmpty(registerAttempt)

}
