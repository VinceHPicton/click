package httpserver

import (
	"fmt"
	"net/http"
	"strings"

	"vincehpicton/click/internal/db/factory"
)

func (ts *HandlerSuite) TestDecodeJSON_EmptyBodyRejected() {
	w := ts.postJSON(ts.routeURL(authAttemptStartLoginRouteName), nil)
	ts.Equal(http.StatusBadRequest, w.Code)

	attempts, err := ts.server.Queries.GetAuthAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(attempts))
}

// A second JSON value after the first would otherwise be silently dropped,
// which lets a caller smuggle a body past anything that inspects only the
// first object. decodeJSON checks decoder.More() for exactly this.
func (ts *HandlerSuite) TestDecodeJSON_TrailingContentRejected() {
	body := []byte(`{"mobile":"+447840195452"}{"mobile":"+447840195453"}`)

	w := ts.postJSON(ts.routeURL(authAttemptStartLoginRouteName), body)
	ts.Equal(http.StatusBadRequest, w.Code)

	attempts, err := ts.server.Queries.GetAuthAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(attempts))
}

func (ts *HandlerSuite) TestDecodeJSON_OversizedBodyRejected() {
	oversized := []byte(fmt.Sprintf(`{"mobile":"%s"}`, strings.Repeat("1", maxBodyBytes+1)))

	w := ts.postJSON(ts.routeURL(authAttemptStartLoginRouteName), oversized)
	ts.Equal(http.StatusBadRequest, w.Code)

	attempts, err := ts.server.Queries.GetAuthAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(attempts))
}

func (ts *HandlerSuite) TestDecodeJSON_BodyJustUnderLimitAccepted() {
	w := ts.postJSON(ts.routeURL(authAttemptStartLoginRouteName), []byte(`{"mobile":"+447840195452"}`))
	ts.Equal(http.StatusOK, w.Code)
}

func (ts *HandlerSuite) TestDecodeJSON_UnknownFieldRejectedOnConfirmLogin() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, user.Mobile)
	ts.Require().NoError(err)

	body := []byte(fmt.Sprintf(
		`{"id":%q,"oneTimeCode":%d,"admin":true}`,
		authAttempt.ID, authAttempt.OneTimeCode,
	))

	w := ts.postJSON(ts.routeURL(authAttemptConfirmLoginRouteName), body)
	ts.Equal(http.StatusBadRequest, w.Code)

	ts.False(ts.onlyAuthAttempt().UsedAt.Valid, "a rejected body must not consume the attempt")

	refreshTokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(refreshTokens))
}

func (ts *HandlerSuite) TestDecodeJSON_UnknownFieldRejectedOnConfirmCreateUser() {
	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	body := []byte(fmt.Sprintf(
		`{"id":%q,"oneTimeCode":%d,"admin":true}`,
		authAttempt.ID, authAttempt.OneTimeCode,
	))

	w := ts.postJSON(ts.routeURL(authAttemptConfirmCreateUserRouteName), body)
	ts.Equal(http.StatusBadRequest, w.Code)

	ts.False(ts.onlyAuthAttempt().UsedAt.Valid)

	users, err := ts.server.Queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(0, len(users))
}

func (ts *HandlerSuite) TestDecodeJSON_UnknownFieldRejectedOnRefresh() {
	w := ts.postJSON(ts.routeURL(refreshRouteName), []byte(`{"refreshToken":"x","admin":true}`))
	ts.Equal(http.StatusBadRequest, w.Code)
}

func (ts *HandlerSuite) TestDecodeJSON_UnknownFieldRejectedOnLogout() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	body := []byte(fmt.Sprintf(`{"refreshToken":%q,"admin":true}`, refreshToken))

	w := ts.postJSONAs(ts.routeURL(logoutRouteName), body, user.ID.String())
	ts.Equal(http.StatusBadRequest, w.Code)

	tokenRows, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Len(tokenRows, 1)
	ts.False(tokenRows[0].RevokedAt.Valid, "a rejected body must not revoke the token")
}
