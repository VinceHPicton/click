package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"
	"vincehpicton/click/internal/db/factory"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"
)

func (ts *HandlerSuite) TestLogout() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	w := ts.callLogout(refreshToken, user.ID.String())
	ts.Equal(http.StatusOK, w.Code)

	tokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(1, len(tokens))
	token := tokens[0]

	ts.Equal(true, token.RevokedAt.Valid)
}

func (ts *HandlerSuite) TestLogout_BannedUser() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	err = ts.server.Queries.BanUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	w := ts.callLogout(refreshToken, user.ID.String())
	ts.Equal(http.StatusOK, w.Code)

	tokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(1, len(tokens))
	token := tokens[0]

	ts.Equal(true, token.RevokedAt.Valid)
}

func (ts *HandlerSuite) TestLogout_InvalidRefreshToken() {
	w := ts.callLogout("invalid-token", "")
	ts.Equal(http.StatusUnauthorized, w.Code)
}

func (ts *HandlerSuite) TestLogout_ExpiredRefresh() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken, dbToken, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	err = ts.server.Queries.SetRefreshTokenExpiryByID(ts.ctx, sqlc.SetRefreshTokenExpiryByIDParams{
		ID:        dbToken.ID,
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	})
	ts.Require().NoError(err)

	w := ts.callLogout(refreshToken, user.ID.String())
	ts.Equal(http.StatusBadRequest, w.Code)

	tokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(1, len(tokens))
	token := tokens[0]

	ts.Equal(false, token.RevokedAt.Valid)
}

func (ts *HandlerSuite) callLogout(refreshToken string, userID string) *httptest.ResponseRecorder {

	token, err := ts.server.TokenManager.GenerateAccessToken(userID)
	ts.Require().NoError(err)

	body, err := json.Marshal(logoutRequest{
		RefreshToken: refreshToken,
	})
	ts.Require().NoError(err)

	url, err := ts.server.Router.Get(logoutRouteName).URL()
	ts.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		url.String(),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	ts.server.Router.ServeHTTP(w, req)

	return w
}

// TestLogout_InvalidRefreshToken covers the 401 from the auth middleware. This
// covers the case past it: a caller who is authenticated but presents a refresh
// token that was never issued.
func (ts *HandlerSuite) TestLogout_AuthenticatedWithUnknownRefreshToken() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	unknownToken, err := tokens.GenerateRefreshToken()
	ts.Require().NoError(err)

	w := ts.callLogout(unknownToken, user.ID.String())
	ts.Equal(http.StatusBadRequest, w.Code)
}

// Logging out twice with the same token is a client bug, not a server error.
func (ts *HandlerSuite) TestLogout_SecondLogoutRejected() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	first := ts.callLogout(refreshToken, user.ID.String())
	ts.Require().Equal(http.StatusOK, first.Code)

	second := ts.callLogout(refreshToken, user.ID.String())
	ts.Equal(http.StatusBadRequest, second.Code)
}

// SECURITY GAP, pinned so it is not mistaken for intended behaviour.
//
// logoutHandler never reads AuthFromContext, and service.Logout takes only the
// token string, so any authenticated caller can revoke any refresh token whose
// value they know. The access token's subject is ignored entirely. Logout
// should be scoped to the calling user; when it is, this test should expect a
// rejection and the victim's token to survive.
func (ts *HandlerSuite) TestLogout_RevokesAnotherUsersRefreshToken() {
	attacker, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	victim, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	victimToken, victimDBToken, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, victim.ID)
	ts.Require().NoError(err)

	w := ts.callLogout(victimToken, attacker.ID.String())
	ts.Equal(http.StatusOK, w.Code)

	revoked, err := ts.server.Queries.GetTokenByHash(ts.ctx, victimDBToken.TokenHash)
	ts.Require().NoError(err)
	ts.True(revoked.RevokedAt.Valid, "victim's token was revoked by an unrelated caller")
}

func (ts *HandlerSuite) TestLogout_RequiresAuthentication() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	body, err := json.Marshal(logoutRequest{RefreshToken: refreshToken})
	ts.Require().NoError(err)

	w := ts.postJSON(ts.routeURL(logoutRouteName), body)
	ts.Equal(http.StatusUnauthorized, w.Code)

	tokenRows, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Len(tokenRows, 1)
	ts.False(tokenRows[0].RevokedAt.Valid)
}

// Logout revokes the refresh token but cannot revoke the access token, which
// has no deny list. The access token stays usable until it expires.
func (ts *HandlerSuite) TestLogout_LeavesAccessTokenUsable() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	loggedOut := ts.callLogout(refreshToken, user.ID.String())
	ts.Require().Equal(http.StatusOK, loggedOut.Code)

	stillAuthorised := ts.sendRequest(http.MethodDelete, ts.routeURL(userDeleteRouteName), nil, map[string]string{
		"Authorization": "Bearer " + ts.accessTokenFor(user.ID.String()),
	})
	ts.Equal(http.StatusOK, stillAuthorised.Code)
}
