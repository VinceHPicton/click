package httpserver

import (
	"net/http"
)

func (ts *HandlerSuite) TestDBPing() {
	w := ts.sendRequest(http.MethodGet, ts.routeURL(dbPingRouteName), nil, nil)

	ts.Equal(http.StatusOK, w.Code)
	ts.Equal("DB pinged", w.Body.String())
}

// /dbping has no .Methods matcher, so middlewareExample is the only thing
// restricting it to GET.
func (ts *HandlerSuite) TestDBPing_RejectsNonGET() {
	for _, method := range []string{http.MethodPost, http.MethodDelete, http.MethodPut} {
		w := ts.sendRequest(method, ts.routeURL(dbPingRouteName), nil, nil)
		ts.Equal(http.StatusMethodNotAllowed, w.Code, method)
	}
}

// UNAUTHENTICATED DIAGNOSTIC ENDPOINT, pinned so the exposure is on record.
//
// /dbping sits outside the auth middleware and, on failure, writes the raw
// driver error to the response with a 200 status. It should either require auth
// or return a bare status code. Adjust this test when it is locked down.
func (ts *HandlerSuite) TestDBPing_NeedsNoAuthentication() {
	w := ts.sendRequest(http.MethodGet, ts.routeURL(dbPingRouteName), nil, nil)
	ts.Equal(http.StatusOK, w.Code)
}
