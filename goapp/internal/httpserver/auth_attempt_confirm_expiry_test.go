package httpserver

import (
	"net/http"
	"time"

	"vincehpicton/click/internal/db/factory"
)

// GetValidAuthAttempt only accepts attempts created within the last two
// minutes. Nothing else enforces that window, so both confirm endpoints are
// checked against an aged attempt.
func (ts *HandlerSuite) TestConfirmLogin_ExpiredAttemptRejected() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, user.Mobile)
	ts.Require().NoError(err)

	ts.setAuthAttemptCreatedAt(authAttempt.ID, authAttempt.CreatedAt.Time.Add(-130*time.Second))

	w := ts.callConfirmLogin(confirmLoginRequest{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode,
	})
	ts.Equal(http.StatusUnauthorized, w.Code)

	refreshTokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(refreshTokens))
	ts.False(ts.onlyAuthAttempt().UsedAt.Valid, "an expired attempt is never reached, so it is not consumed")
}

func (ts *HandlerSuite) TestConfirmLogin_AttemptInsideWindowAccepted() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, user.Mobile)
	ts.Require().NoError(err)

	ts.setAuthAttemptCreatedAt(authAttempt.ID, authAttempt.CreatedAt.Time.Add(-110*time.Second))

	w := ts.callConfirmLogin(confirmLoginRequest{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode,
	})
	ts.Equal(http.StatusOK, w.Code)
}

func (ts *HandlerSuite) TestConfirmCreateUser_ExpiredAttemptRejected() {
	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	ts.setAuthAttemptCreatedAt(authAttempt.ID, authAttempt.CreatedAt.Time.Add(-130*time.Second))

	w := ts.callConfirmCreateUser(authAttempt.ID, authAttempt.OneTimeCode)
	ts.Equal(http.StatusUnauthorized, w.Code)

	users, err := ts.server.Queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(users))
}

func (ts *HandlerSuite) TestConfirmCreateUser_AttemptInsideWindowAccepted() {
	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	ts.setAuthAttemptCreatedAt(authAttempt.ID, authAttempt.CreatedAt.Time.Add(-110*time.Second))

	w := ts.callConfirmCreateUser(authAttempt.ID, authAttempt.OneTimeCode)
	ts.Equal(http.StatusOK, w.Code)
}
