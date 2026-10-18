package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCORSMiddleware_SetsHeadersAndCallsNext(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	recorder := httptest.NewRecorder()

	CORSMiddleware()(next).ServeHTTP(recorder, request)

	require.True(t, nextCalled)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "http://localhost:8081", recorder.Header().Get("Access-Control-Allow-Origin"))
	require.Contains(t, recorder.Header().Get("Access-Control-Allow-Methods"), http.MethodPost)
	require.Contains(t, recorder.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	require.Contains(t, recorder.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
}

// A preflight must be answered by the middleware itself. Passing OPTIONS down
// the chain would hand it to an auth check that a browser preflight cannot
// satisfy, since preflights carry no Authorization header.
func TestCORSMiddleware_PreflightShortCircuits(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	})

	request := httptest.NewRequest(http.MethodOptions, "/auth/logout", nil)
	recorder := httptest.NewRecorder()

	CORSMiddleware()(next).ServeHTTP(recorder, request)

	require.False(t, nextCalled, "preflight must not reach the wrapped handler")
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "http://localhost:8081", recorder.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_PreflightShortCircuitsBeforeAuth(t *testing.T) {
	handlerCalled := false
	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})

	tokenMgr := newTestTokenManager(t)
	chained := CORSMiddleware()(AuthMiddleware(tokenMgr)(handler))

	request := httptest.NewRequest(http.MethodOptions, "/auth/logout", nil)
	recorder := httptest.NewRecorder()

	chained.ServeHTTP(recorder, request)

	require.False(t, handlerCalled)
	require.Equal(t, http.StatusOK, recorder.Code, "preflight must not be rejected by auth")
}

// KNOWN DEFECT, pinned so it is not mistaken for working behaviour.
//
// Every route is registered with an explicit .Methods(...) matcher that omits
// OPTIONS. gorilla/mux only builds the middleware chain when the match
// succeeds (Router.Match skips it when MatchErr is set), so a preflight fails
// method matching, never reaches CORSMiddleware, and comes back as a bare 405
// with no CORS headers. Browsers will refuse the follow-up request.
//
// The middleware is correct; the routing is not. Fixing it means allowing
// OPTIONS on these routes, or installing a MethodNotAllowedHandler that
// answers preflights. Update this test to expect 200 when that lands.
func (ts *HandlerSuite) TestCORS_PreflightNeverReachesMiddleware() {
	preflightTargets := []string{
		ts.routeURL(authAttemptStartLoginRouteName),
		ts.routeURL(authAttemptConfirmLoginRouteName),
		ts.routeURL(refreshRouteName),
		ts.routeURL(logoutRouteName),
		ts.routeURL(userDeleteRouteName),
	}

	for _, target := range preflightTargets {
		w := ts.sendRequest(http.MethodOptions, target, nil, map[string]string{
			"Origin":                         "http://localhost:8081",
			"Access-Control-Request-Method":  http.MethodPost,
			"Access-Control-Request-Headers": "Content-Type, Authorization",
		})

		ts.Equal(http.StatusMethodNotAllowed, w.Code, "preflight to %s", target)
		ts.Empty(w.Header().Get("Access-Control-Allow-Origin"), "preflight to %s", target)
	}
}
