package httpserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"

	"github.com/stretchr/testify/require"
)

func (ts *HandlerSuite) routeURL(name string) string {
	url, err := ts.server.Router.Get(name).URL()
	ts.Require().NoError(err)
	return url.String()
}

func (ts *HandlerSuite) accessTokenFor(userID string) string {
	token, err := ts.server.TokenManager.GenerateAccessToken(userID)
	ts.Require().NoError(err)
	return token
}

func (ts *HandlerSuite) sendRequest(method, target string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	w := httptest.NewRecorder()
	ts.server.Router.ServeHTTP(w, req)
	return w
}

func (ts *HandlerSuite) postJSON(target string, body []byte) *httptest.ResponseRecorder {
	return ts.sendRequest(http.MethodPost, target, body, map[string]string{
		"Content-Type": "application/json",
	})
}

func (ts *HandlerSuite) postJSONAs(target string, body []byte, userID string) *httptest.ResponseRecorder {
	return ts.sendRequest(http.MethodPost, target, body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + ts.accessTokenFor(userID),
	})
}

func (ts *HandlerSuite) onlyAuthAttempt() sqlc.AppAuthAttempt {
	attempts, err := ts.server.Queries.GetAuthAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Len(attempts, 1)
	return attempts[0]
}

func newTestTokenManager(t *testing.T) *tokens.Manager {
	tokenMgr, err := tokens.New("helper-test-secret")
	require.NoError(t, err)
	return tokenMgr
}
