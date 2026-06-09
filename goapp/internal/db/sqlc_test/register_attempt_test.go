package sqlc_test

import "vincehpicton/click/internal/db/sqlc"

func (ts *DatabaseSuite) TestAuthAttemptCreate() {
	authAttempt, err := ts.queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)
	ts.NotEmpty(authAttempt)

	authAttempts, err := ts.queries.GetAuthAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(authAttempts))
	ts.Equal(authAttempt.ID, authAttempts[0].ID)
}

func (ts *DatabaseSuite) TestAuthAttemptConfirm() {
	authAttempt, err := ts.queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	AuthAttemptConfirmParams := sqlc.AuthAttemptConfirmParams{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode,
	}

	newUserID, err := ts.queries.AuthAttemptConfirm(ts.ctx, AuthAttemptConfirmParams)
	ts.Require().NoError(err)
	ts.NotEmpty(newUserID)
	ts.NotEqual(newUserID, authAttempt.ID)

	users, err := ts.queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(users))
	ts.NotEmpty(users)
}

func (ts *DatabaseSuite) TestAuthAttemptConfirm_Fail() {
	authAttempt, err := ts.queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	AuthAttemptConfirmParams := sqlc.AuthAttemptConfirmParams{
		OneTimeCode: authAttempt.OneTimeCode,
	}

	_, err = ts.queries.AuthAttemptConfirm(ts.ctx, AuthAttemptConfirmParams)
	ts.Require().Error(err)

	users, err := ts.queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(users))
	ts.Empty(users)
}
