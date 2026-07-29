package httpserver // adjust to your actual package name

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"vincehpicton/click/internal/tokens"
)

// buildAuthHeader is a helper type for table cases that need a live token
// generated from a specific tokens.Manager (e.g. valid token, wrong-signer
// token). Cases that don't need a real token (missing header, malformed
// header, garbage string) just return a static header directly.
type buildAuthHeaderFunc func(t *testing.T, tokenMgr *tokens.Manager) string

func TestAuthMiddleware(t *testing.T) {
	const secret = "super-secret-code"
	const userID = "fake-uuid"

	tokenMgr, err := tokens.New(secret)
	require.NoError(t, err)

	// A second manager with a different signing secret, used to simulate
	// a token that is well-formed but was never issued by this service.
	otherTokenMgr, err := tokens.New("a-different-secret-code")
	require.NoError(t, err)

	// A manager whose clock is pinned an hour in the past. Tokens it mints
	// have already expired by the time the (real-clock) tokenMgr parses
	// them, letting us exercise the expiry path without hand-crafting a
	// JWT or reaching into the tokens package's internals.
	expiredTokenMgr, err := tokens.New(secret, tokens.WithNowFunc(func() time.Time {
		return time.Now().Add(-1 * time.Hour)
	}))
	require.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string              // static header value; used if buildHeader is nil
		buildHeader    buildAuthHeaderFunc // dynamic header value; takes precedence over authHeader
		wantStatus     int
		wantNextCalled bool
		wantUserID     string // only checked when wantNextCalled is true
	}{
		{
			name: "valid Bearer token (capitalized)",
			buildHeader: func(t *testing.T, tm *tokens.Manager) string {
				tok, err := tm.GenerateAccessToken(userID)
				require.NoError(t, err)
				return "Bearer " + tok
			},
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
			wantUserID:     userID,
		},
		{
			name: "valid bearer token (lowercase)",
			buildHeader: func(t *testing.T, tm *tokens.Manager) string {
				tok, err := tm.GenerateAccessToken(userID)
				require.NoError(t, err)
				return "bearer " + tok
			},
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
			wantUserID:     userID,
		},
		{
			name:           "missing Authorization header",
			authHeader:     "",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "header with no space (unparsable scheme/token)",
			authHeader:     "BearerTokenGlued",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "unsupported scheme",
			authHeader:     "Basic dXNlcjpwYXNz",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "scheme wrong case (BEARER)",
			authHeader:     "BEARER sometoken",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "Bearer with empty token",
			authHeader:     "Bearer ",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "Bearer with garbage/malformed token",
			authHeader:     "Bearer not-a-real-jwt",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name: "token signed by a different manager/secret",
			buildHeader: func(t *testing.T, _ *tokens.Manager) string {
				tok, err := otherTokenMgr.GenerateAccessToken(userID)
				require.NoError(t, err)
				return "Bearer " + tok
			},
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name: "expired token",
			buildHeader: func(t *testing.T, _ *tokens.Manager) string {
				tok, err := expiredTokenMgr.GenerateAccessToken(userID)
				require.NoError(t, err)
				return "Bearer " + tok
			},
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "extra whitespace between scheme and token",
			authHeader:     "Bearer  extra-space-token", // note double space
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			var capturedReq *http.Request
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				capturedReq = r
				w.WriteHeader(http.StatusOK)
			})

			request := httptest.NewRequest(http.MethodGet, "/users", nil)

			header := tt.authHeader
			if tt.buildHeader != nil {
				header = tt.buildHeader(t, tokenMgr)
			}
			if header != "" {
				request.Header.Set("Authorization", header)
			}

			recorder := httptest.NewRecorder()
			AuthMiddleware(tokenMgr)(next).ServeHTTP(recorder, request)

			require.Equal(t, tt.wantStatus, recorder.Code)
			require.Equal(t, tt.wantNextCalled, nextCalled, "next handler invocation mismatch")

			if tt.wantNextCalled {
				authCtx, ok := AuthFromContext(capturedReq.Context())
				require.True(t, ok, "expected auth context to be present")
				require.Equal(t, tt.wantUserID, authCtx.UserID)
			}
		})
	}
}

// TestAuthFromContext_NoValue guards the helper itself: calling it on a
// context that was never populated by the middleware should report ok=false
// rather than panicking or returning a zero-value AuthContext silently.
func TestAuthFromContext_NoValue(t *testing.T) {
	_, ok := AuthFromContext(context.Background())
	require.False(t, ok)
}

// My tests for my understanding of middlware testing
func TestAuthMiddleware_HappyPath(t *testing.T) {
	tokenMgr, err := tokens.New("super-secret-code")
	if err != nil {
		t.Fatal(err)
	}
	const userID = "fake uuid"
	token, err := tokenMgr.GenerateAccessToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	cases := []string{"Bearer", "bearer"}
	for _, scheme := range cases {
		t.Run(scheme, func(t *testing.T) {
			var capturedReq *http.Request
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedReq = r
				w.WriteHeader(http.StatusOK)
			})

			request := httptest.NewRequest(http.MethodGet, "/users", nil)
			request.Header.Set("Authorization", scheme+" "+token)
			recorder := httptest.NewRecorder()

			AuthMiddleware(tokenMgr)(next).ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			authCtx, ok := AuthFromContext(capturedReq.Context())
			require.True(t, ok, "expected auth context")
			require.Equal(t, userID, authCtx.UserID)
		})
	}
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	request := httptest.NewRequest(http.MethodGet, "/users", nil)
	recorder := httptest.NewRecorder()

	tokenMgr, err := tokens.New("super-secret-code")
	if err != nil {
		t.Fatal(err)
	}

	AuthMiddleware(tokenMgr)(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}
