package sqlc_test

import "vincehpicton/click/internal/db/sqlc"

func (ts *DatabaseSuite) TestRegisterAttemptCreate() {
	registerAttempt, err := ts.queries.RegisterAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)
	ts.NotEmpty(registerAttempt)

	registerAttempts, err := ts.queries.GetRegisterAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(registerAttempts))
	ts.Equal(registerAttempt.ID, registerAttempts[0].ID)
}

func (ts *DatabaseSuite) TestRegisterAttemptConfirm() {
	registerAttempt, err := ts.queries.RegisterAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	RegisterAttemptConfirmParams := sqlc.RegisterAttemptConfirmParams{
		ID:          registerAttempt.ID,
		OneTimeCode: registerAttempt.OneTimeCode,
	}

	newUserID, err := ts.queries.RegisterAttemptConfirm(ts.ctx, RegisterAttemptConfirmParams)
	ts.Require().NoError(err)
	ts.NotEmpty(newUserID)
	ts.NotEqual(newUserID, registerAttempt.ID)

	users, err := ts.queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(users))
	ts.NotEmpty(users)
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
