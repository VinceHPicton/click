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

func (ts *DatabaseSuite) TestAuthAttemptCreateUser() {
	authAttempt, err := ts.queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	AuthAttemptCreateUserParams := sqlc.AuthAttemptCreateUserParams{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode,
	}

	newUserID, err := ts.queries.AuthAttemptCreateUser(ts.ctx, AuthAttemptCreateUserParams)
	ts.Require().NoError(err)
	ts.NotEmpty(newUserID)
	ts.NotEqual(newUserID, authAttempt.ID)

	users, err := ts.queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(users))
	ts.NotEmpty(users)
}

func (ts *DatabaseSuite) TestAuthAttemptCreateUser_Fail() {
	authAttempt, err := ts.queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	AuthAttemptCreateUserParams := sqlc.AuthAttemptCreateUserParams{
		OneTimeCode: authAttempt.OneTimeCode,
	}

	_, err = ts.queries.AuthAttemptCreateUser(ts.ctx, AuthAttemptCreateUserParams)
	ts.Require().Error(err)

	users, err := ts.queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(users))
	ts.Empty(users)
}
