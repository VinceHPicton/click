package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"vincehpicton/click/internal/db/factory"
)

func (ts *HandlerSuite) TestConfirmLogin() {
	var err error

	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, user.Mobile)
	ts.Require().NoError(err)

	req := confirmLoginRequest{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode,
	}

	w := ts.callConfirmLogin(req)

	ts.Equal(http.StatusOK, w.Code)

	refreshTokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(refreshTokens))

	resp := confirmLoginResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	ts.Require().NoError(err)

	ts.Require().NotEmpty(resp.UserID)
	ts.Require().NotEmpty(resp.AccessToken)
	ts.Require().NotEmpty(resp.RefreshToken)
}

func (ts *HandlerSuite) TestConfirmLogin_ThenRefresh() {
	var err error

	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, user.Mobile)
	ts.Require().NoError(err)

	req := confirmLoginRequest{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode,
	}

	w := ts.callConfirmLogin(req)

	ts.Equal(http.StatusOK, w.Code)

	refreshTokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(refreshTokens))

	resp := confirmLoginResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	ts.Require().NoError(err)

	refreshW := ts.callRefresh(resp.RefreshToken)

	ts.Equal(http.StatusOK, refreshW.Code)

	refreshResp := refreshResponse{}
	err = json.Unmarshal(refreshW.Body.Bytes(), &refreshResp)
	ts.Require().NoError(err)

	ts.NotEmpty(refreshResp.AccessToken)
	ts.NotEmpty(refreshResp.RefreshToken)

	refreshTokens, err = ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	// 2 because refresh token doesnt hard delete the old one
	ts.Equal(2, len(refreshTokens))
}

func (ts *HandlerSuite) TestConfirmLogin_UserDoesntExist() {
	var err error
	const randomMobile = "+447840195452"

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, randomMobile)
	ts.Require().NoError(err)

	req := confirmLoginRequest{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode,
	}

	w := ts.callConfirmLogin(req)

	ts.Equal(http.StatusNotFound, w.Code)

	refreshTokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(0, len(refreshTokens))
}

func (ts *HandlerSuite) TestConfirmLogin_UserBanned() {
	var err error

	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	err = ts.server.Queries.BanUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, user.Mobile)
	ts.Require().NoError(err)

	req := confirmLoginRequest{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode,
	}

	w := ts.callConfirmLogin(req)

	// TODO: banned users are not found, this might be fine; ie can show user "Your account is banned, or doesn't exist"
	ts.Equal(http.StatusNotFound, w.Code)

	refreshTokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(0, len(refreshTokens))
}

func (ts *HandlerSuite) TestConfirmLogin_UserDeleted() {
	var err error

	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	err = ts.server.Queries.SoftDeleteUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, user.Mobile)
	ts.Require().NoError(err)

	req := confirmLoginRequest{
		ID:          authAttempt.ID,
		OneTimeCode: authAttempt.OneTimeCode,
	}

	w := ts.callConfirmLogin(req)

	ts.Equal(http.StatusNotFound, w.Code)

	refreshTokens, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(0, len(refreshTokens))
}

func (ts *HandlerSuite) callConfirmLogin(reqStruct confirmLoginRequest) *httptest.ResponseRecorder {
	body, err := json.Marshal(reqStruct)
	ts.Require().NoError(err)

	url, err := ts.server.Router.Get(authAttemptConfirmLoginRouteName).URL()
	ts.Require().NoError(err)

	httpReq := httptest.NewRequest(
		http.MethodPost,
		url.String(),
		bytes.NewReader(body),
	)
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler := ts.server.authAttemptConfirmLoginHandler()
	handler(w, httpReq)

	return w
}
