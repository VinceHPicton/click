package httpserver

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"

	"github.com/google/uuid"
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

// The offset is applied to the created_at Postgres assigned, so the validity
// window assertions do not depend on the Go process and the database container
// agreeing on the current time.
func (ts *HandlerSuite) setAuthAttemptCreatedAt(id uuid.UUID, createdAt time.Time) {
	err := ts.server.Queries.SetAuthAttemptCreatedAt(ts.ctx, sqlc.SetAuthAttemptCreatedAtParams{
		ID:        id,
		CreatedAt: sql.NullTime{Time: createdAt, Valid: true},
	})
	ts.Require().NoError(err)
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
