package httpserver

import (
	"net/http"

	"vincehpicton/click/internal/db/factory"

	"github.com/google/uuid"
)

func (ts *HandlerSuite) TestConfirmLogin_WrongCodeRejected() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, user.Mobile)
	ts.Require().NoError(err)

	w := ts.callConfirmLogin(confirmLoginRequest{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode + 1,
	})
	ts.Equal(http.StatusUnauthorized, w.Code)

	refreshTokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(refreshTokens))
}

// A wrong code burns the attempt. Without this the code is brute-forceable
// within the two minute window, so the consume-before-compare ordering in
// service.getAndConsumeAttempt must not be reordered.
func (ts *HandlerSuite) TestConfirmLogin_WrongCodeStillConsumesAttempt() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, user.Mobile)
	ts.Require().NoError(err)

	wrong := ts.callConfirmLogin(confirmLoginRequest{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode + 1,
	})
	ts.Require().Equal(http.StatusUnauthorized, wrong.Code)

	ts.True(ts.onlyAuthAttempt().UsedAt.Valid, "wrong code must still consume the attempt")

	retry := ts.callConfirmLogin(confirmLoginRequest{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode,
	})
	ts.Equal(http.StatusUnauthorized, retry.Code, "the correct code must not work after a wrong guess")

	refreshTokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(refreshTokens))
}

// writeError deliberately collapses ErrInvalidAttempt and ErrInvalidCode onto
// one response so a caller cannot use the status or body to work out whether
// the attempt ID was real.
func (ts *HandlerSuite) TestConfirmLogin_WrongCodeIndistinguishableFromUnknownAttempt() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, user.Mobile)
	ts.Require().NoError(err)

	wrongCode := ts.callConfirmLogin(confirmLoginRequest{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode + 1,
	})
	unknownAttempt := ts.callConfirmLogin(confirmLoginRequest{
		ID:          uuid.New(),
		OneTimeCode: authAttempt.OneTimeCode,
	})

	ts.Equal(http.StatusUnauthorized, wrongCode.Code)
	ts.Equal(wrongCode.Code, unknownAttempt.Code)
	ts.Equal(wrongCode.Body.String(), unknownAttempt.Body.String())
}

func (ts *HandlerSuite) TestConfirmCreateUser_WrongCodeRejected() {
	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	w := ts.callConfirmCreateUser(authAttempt.ID, authAttempt.OneTimeCode+1)
	ts.Equal(http.StatusUnauthorized, w.Code)

	users, err := ts.server.Queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(users), "a wrong code must not create an account")
}

func (ts *HandlerSuite) TestConfirmCreateUser_WrongCodeStillConsumesAttempt() {
	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	wrong := ts.callConfirmCreateUser(authAttempt.ID, authAttempt.OneTimeCode+1)
	ts.Require().Equal(http.StatusUnauthorized, wrong.Code)

	ts.True(ts.onlyAuthAttempt().UsedAt.Valid)

	retry := ts.callConfirmCreateUser(authAttempt.ID, authAttempt.OneTimeCode)
	ts.Equal(http.StatusUnauthorized, retry.Code)

	users, err := ts.server.Queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(users))
}

func (ts *HandlerSuite) TestConfirmCreateUser_WrongCodeIndistinguishableFromUnknownAttempt() {
	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	wrongCode := ts.callConfirmCreateUser(authAttempt.ID, authAttempt.OneTimeCode+1)
	unknownAttempt := ts.callConfirmCreateUser(uuid.New(), authAttempt.OneTimeCode)

	ts.Equal(http.StatusUnauthorized, wrongCode.Code)
	ts.Equal(wrongCode.Code, unknownAttempt.Code)
	ts.Equal(wrongCode.Body.String(), unknownAttempt.Body.String())
}
