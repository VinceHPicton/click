package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"vincehpicton/click/internal/db/factory"
)

func (ts *HandlerSuite) TestLogout() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	w := ts.callLogout(refreshToken)
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

	w := ts.callLogout(refreshToken)
	ts.Equal(http.StatusOK, w.Code)

	tokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(1, len(tokens))
	token := tokens[0]

	ts.Equal(true, token.RevokedAt.Valid)
}

func (ts *HandlerSuite) TestLogout_InvalidRefreshToken() {
	w := ts.callLogout("invalid-token")
	ts.Equal(http.StatusBadRequest, w.Code)
}

// func (ts *HandlerSuite) TestLogout_ExpiredRefresh() {
// 	var err error
// 	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
// 	ts.Require().NoError(err)

// 	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
// 	ts.Require().NoError(err)

// 	err = ts.server.Queries.SetRefreshTokenExpiryByID(ts.ctx, sqlc.SetRefreshTokenExpiryByIDParams{
// 		ExpiresAt: time.Now().Add(-24 * time.Hour),
// 	})
// 	ts.Require().NoError(err)

// 	w := ts.callLogout(refreshToken)
// 	ts.Equal(http.StatusBadRequest, w.Code)

// 	tokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
// 	ts.Require().NoError(err)
// 	ts.Equal(1, len(tokens))
// 	token := tokens[0]

// 	ts.Equal(false, token.RevokedAt.Valid)
// }

func (ts *HandlerSuite) callLogout(refreshToken string) *httptest.ResponseRecorder {
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

	w := httptest.NewRecorder()

	handler := ts.server.logoutHandler()
	handler(w, req)

	return w
}
