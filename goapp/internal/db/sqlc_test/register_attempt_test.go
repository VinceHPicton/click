package sqlc_test

import "vincehpicton/click/internal/db/sqlc"

func (ts *DatabaseSuite) TestRegisterAttemptCreate() {
	registerAttempt, err := ts.queries.RegisterAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)
	ts.NotEmpty(registerAttempt)
}

func (ts *DatabaseSuite) TestRegisterAttemptConfirm() {
	registerAttempt, err := ts.queries.RegisterAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	RegisterAttemptConfirmParams := sqlc.RegisterAttemptConfirmParams{
		ID:          registerAttempt.ID,
		OneTimeCode: registerAttempt.OneTimeCode,
	}

	newUser, err := ts.queries.RegisterAttemptConfirm(ts.ctx, RegisterAttemptConfirmParams)
	ts.Require().NoError(err)
	ts.NotEmpty(newUser)
}

func (ts *DatabaseSuite) TestRegisterAttemptConfirm_Fail() {
	registerAttempt, err := ts.queries.RegisterAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	RegisterAttemptConfirmParams := sqlc.RegisterAttemptConfirmParams{
		OneTimeCode: registerAttempt.OneTimeCode,
	}

	_, err = ts.queries.RegisterAttemptConfirm(ts.ctx, RegisterAttemptConfirmParams)
	ts.Require().Error(err)

	users, err := ts.queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(users))
	ts.Empty(users)
}
