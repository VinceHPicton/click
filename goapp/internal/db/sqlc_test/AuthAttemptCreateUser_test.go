package sqlc_test

// -- name: AuthAttemptCreateUser :one

// WITH consumed_attempt AS (
  
//     UPDATE app.auth_attempts AS ra
//     SET used_at = NOW()
//     WHERE ra.id = sqlc.arg(id)
//       AND ra.one_time_code = sqlc.arg(one_time_code)
//       AND ra.used_at IS NULL
//       AND ra.created_at >= NOW() - INTERVAL '2 minutes'
//     RETURNING id, mobile
// ),
// new_user AS (
//     INSERT INTO app.users (mobile)
//     SELECT mobile
//     FROM consumed_attempt
//     RETURNING id
// )
// SELECT id AS user_id
// FROM new_user;

// func (ts *DatabaseSuite) TestAuthAttemptCreate() {
// 	authAttempt, err := ts.queries.AuthAttemptCreate(ts.ctx, "+447840195455")
// 	ts.Require().NoError(err)
// 	ts.NotEmpty(authAttempt)

// 	authAttempts, err := ts.queries.GetAuthAttempts(ts.ctx)
// 	ts.Require().NoError(err)
// 	ts.Require().Equal(1, len(authAttempts))
// 	ts.Equal(authAttempt.ID, authAttempts[0].ID)
// }

// func (ts *DatabaseSuite) TestAuthAttemptCreateUser() {
// 	authAttempt, err := ts.queries.AuthAttemptCreate(ts.ctx, "+447840195455")
// 	ts.Require().NoError(err)

// 	AuthAttemptCreateUserParams := sqlc.AuthAttemptCreateUserParams{
// 		ID:          authAttempt.ID,
// 		OneTimeCode: authAttempt.OneTimeCode,
// 	}

// 	newUserID, err := ts.queries.AuthAttemptCreateUser(ts.ctx, AuthAttemptCreateUserParams)
// 	ts.Require().NoError(err)
// 	ts.NotEmpty(newUserID)
// 	ts.NotEqual(newUserID, authAttempt.ID)

// 	users, err := ts.queries.GetAllUsers(ts.ctx)
// 	ts.Require().NoError(err)
// 	ts.Require().Equal(1, len(users))
// 	ts.NotEmpty(users)
// }

// func (ts *DatabaseSuite) TestAuthAttemptCreateUser_FailNoID_AttemptConsumed() {
// 	authAttempt, err := ts.queries.AuthAttemptCreate(ts.ctx, "+447840195455")
// 	ts.Require().NoError(err)

// 	AuthAttemptCreateUserParams := sqlc.AuthAttemptCreateUserParams{
// 		OneTimeCode: authAttempt.OneTimeCode,
// 	}

// 	_, err = ts.queries.AuthAttemptCreateUser(ts.ctx, AuthAttemptCreateUserParams)
// 	ts.Require().Error(err)

// 	users, err := ts.queries.GetAllUsers(ts.ctx)
// 	ts.Require().NoError(err)
// 	ts.Equal(0, len(users))
// 	ts.Empty(users)

// 	queriedAuthAttempt, err := ts.queries.GetAuthAttempt(ts.ctx, authAttempt.ID)
// 	ts.Require().NoError(err)
// 	ts.Equal(true, queriedAuthAttempt.UsedAt.Valid)
// }
