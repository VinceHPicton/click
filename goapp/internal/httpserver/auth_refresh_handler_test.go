package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"vincehpicton/click/internal/db/factory"
	"vincehpicton/click/internal/tokens"
)

func (ts *HandlerSuite) TestRefreshToken() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	w := ts.callRefresh(refreshToken)

	ts.Equal(http.StatusOK, w.Code)

	resp := refreshResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	ts.Require().NoError(err)

	ts.Require().NotEmpty(resp.AccessToken)
	ts.Require().NotEmpty(resp.RefreshToken)
}

func (ts *HandlerSuite) TestRefreshToken_BannedUser() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	err = ts.server.Queries.BanUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	w := ts.callRefresh(refreshToken)

	ts.Equal(http.StatusUnauthorized, w.Code)
}

func (ts *HandlerSuite) TestRefreshToken_SoftDeletedUser() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	err = ts.server.Queries.SoftDeleteUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	w := ts.callRefresh(refreshToken)

	ts.Equal(http.StatusBadRequest, w.Code)

	tokenRows, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(tokenRows))
}

func (ts *HandlerSuite) TestRefreshToken_TokenDoesntExist() {
	var err error
	_, err = factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken, err := tokens.GenerateRefreshToken()
	ts.Require().NoError(err)

	w := ts.callRefresh(refreshToken)

	ts.Equal(http.StatusBadRequest, w.Code)

	tokenRows, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(0, len(tokenRows))
}

func (ts *HandlerSuite) callRefresh(refreshToken string) *httptest.ResponseRecorder {
    body, err := json.Marshal(refreshRequest{
        RefreshToken: refreshToken,
    })
    ts.Require().NoError(err)

    url, err := ts.server.Router.Get(refreshRouteName).URL()
    ts.Require().NoError(err)

    req := httptest.NewRequest(
        http.MethodPost,
        url.String(),
        bytes.NewReader(body),
    )
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()

	handler := ts.server.refreshHandler()
	handler(w, req)

    return w
}