package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"vincehpicton/click/internal/db/factory"
	"vincehpicton/click/internal/db/sqlc"
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

func (ts *HandlerSuite) TestRefreshToken_ExpiredToken() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	expiredToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID, factory.WithExpiry(time.Now().Add(-1*time.Hour)))
	ts.Require().NoError(err)

	w := ts.callRefresh(expiredToken)

	ts.Equal(http.StatusBadRequest, w.Code)
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

func (ts *HandlerSuite) TestRefreshToken_OldTokenRejectedAfterRotation() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	oldRefreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	w := ts.callRefresh(oldRefreshToken)
	ts.Require().Equal(http.StatusOK, w.Code)

	resp := refreshResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	ts.Require().NoError(err)
	ts.Require().NotEqual(oldRefreshToken, resp.RefreshToken)

	replay := ts.callRefresh(oldRefreshToken)
	ts.Equal(http.StatusBadRequest, replay.Code)
}

func (ts *HandlerSuite) TestRefreshToken_UnknownTokenMatchesExpiredResponse() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	expiredToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID, factory.WithExpiry(time.Now().Add(-1*time.Hour)))
	ts.Require().NoError(err)

	wreplay := ts.callRefresh(expiredToken)
	unknown := ts.callRefresh("not-a-real-token")

	ts.Equal(wreplay.Code, unknown.Code)
	ts.Equal(wreplay.Body.String(), unknown.Body.String())
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

	ts.server.Router.ServeHTTP(w, req)

	return w
}

// The refresh endpoint is unauthenticated, so an expired token and a token that
// never existed must be indistinguishable. Otherwise a caller can probe which
// token values were once real.
func (ts *HandlerSuite) TestRefreshToken_ExpiredTokenMatchesUnknownResponse() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken, dbToken, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	err = ts.server.Queries.SetRefreshTokenExpiryByID(ts.ctx, sqlc.SetRefreshTokenExpiryByIDParams{
		ID:        dbToken.ID,
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	})
	ts.Require().NoError(err)

	expired := ts.callRefresh(refreshToken)
	unknown := ts.callRefresh("not-a-real-token")

	ts.Equal(http.StatusBadRequest, expired.Code)
	ts.Equal(expired.Code, unknown.Code)
	ts.Equal(expired.Body.String(), unknown.Body.String())
}

// RefreshSession leans on RotateRefreshToken being a single statement: the
// losing caller's UPDATE re-checks revoked_at IS NULL after the row lock clears,
// matches nothing, and the insert is skipped. Without that, two callers racing
// with one leaked token would both walk away with a live session.
func (ts *HandlerSuite) TestRefreshToken_ConcurrentUseIssuesOneSession() {
	const callers = 4

	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	body, err := json.Marshal(refreshRequest{RefreshToken: refreshToken})
	ts.Require().NoError(err)
	target := ts.routeURL(refreshRouteName)

	recorders := make([]*httptest.ResponseRecorder, callers)
	start := make(chan struct{})
	var wg sync.WaitGroup

	for i := range recorders {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			<-start
			ts.server.Router.ServeHTTP(w, req)
			recorders[i] = w
		}(i)
	}
	close(start)
	wg.Wait()

	accepted := 0
	for _, w := range recorders {
		switch w.Code {
		case http.StatusOK:
			accepted++
		case http.StatusBadRequest:
		default:
			ts.Failf("unexpected refresh status", "got %d: %s", w.Code, w.Body.String())
		}
	}
	ts.Equal(1, accepted, "exactly one concurrent refresh may succeed")

	tokenRows, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Len(tokenRows, 2, "the original plus exactly one replacement")

	live := 0
	for _, row := range tokenRows {
		if !row.RevokedAt.Valid {
			live++
		}
	}
	ts.Equal(1, live, "exactly one live refresh token may remain")
}
