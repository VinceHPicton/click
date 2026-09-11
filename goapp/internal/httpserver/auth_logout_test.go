package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"
	"vincehpicton/click/internal/db/factory"
	"vincehpicton/click/internal/db/sqlc"
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

	ts.server.Routes()
	ts.server.Router.ServeHTTP(w, req)

	return w
}
