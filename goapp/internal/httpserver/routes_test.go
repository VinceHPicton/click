package httpserver

import (
	"net/http"
	"sort"

	"vincehpicton/click/internal/db/factory"

	"github.com/gorilla/mux"
)

type routeSpec struct {
	path    string
	methods []string
}

// Locks the public HTTP surface. A route silently changing path or gaining a
// method is otherwise invisible until a client breaks.
func (ts *HandlerSuite) TestRoutes_NamedRouteInventory() {
	want := map[string]routeSpec{
		dbPingRouteName:                       {path: "/dbping"},
		authAttemptStartLoginRouteName:        {path: "/auth-attempt/start/login", methods: []string{http.MethodPost}},
		authAttemptStartCreateUserRouteName:   {path: "/auth-attempt/start/create-user", methods: []string{http.MethodPost}},
		authAttemptConfirmLoginRouteName:      {path: "/auth-attempt/confirm/login", methods: []string{http.MethodPost}},
		authAttemptConfirmCreateUserRouteName: {path: "/auth-attempt/confirm/create-user", methods: []string{http.MethodPost}},
		refreshRouteName:                      {path: "/auth/refresh", methods: []string{http.MethodPost}},
		logoutRouteName:                       {path: "/auth/logout", methods: []string{http.MethodPost}},
		userDeleteRouteName:                   {path: "/users", methods: []string{http.MethodDelete}},
	}

	got := map[string]routeSpec{}
	unnamed := map[string]routeSpec{}

	err := ts.server.Router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			methods = nil
		}
		sort.Strings(methods)

		spec := routeSpec{path: path, methods: methods}
		if name := route.GetName(); name != "" {
			got[name] = spec
			return nil
		}
		unnamed[path] = spec
		return nil
	})
	ts.Require().NoError(err)

	ts.Equal(want, got)
	ts.Equal(map[string]routeSpec{
		"/users/export": {path: "/users/export", methods: []string{http.MethodGet}},
	}, unnamed)
}

// SUSPECT WIRING, pinned so it is not mistaken for a real endpoint.
//
// /users/export is registered against s.logoutHandler(), so a GET with no body
// fails body decoding rather than exporting anything. It is also wrapped in the
// auth middleware twice: once by the protectedCORS subrouter and again by the
// explicit Chain call. Delete the route or give it a real handler, then delete
// this test.
func (ts *HandlerSuite) TestUsersExportRoute_IsWiredToLogoutHandler() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	w := ts.sendRequest(http.MethodGet, "/users/export", nil, map[string]string{
		"Authorization": "Bearer " + ts.accessTokenFor(user.ID.String()),
	})

	ts.Equal(http.StatusBadRequest, w.Code)
	ts.Contains(w.Body.String(), "malformed request body")
}

func (ts *HandlerSuite) TestRoutes_RejectWrongMethod() {
	cases := []struct {
		name   string
		method string
		target string
	}{
		{name: "refresh via GET", method: http.MethodGet, target: ts.routeURL(refreshRouteName)},
		{name: "logout via GET", method: http.MethodGet, target: ts.routeURL(logoutRouteName)},
		{name: "start login via DELETE", method: http.MethodDelete, target: ts.routeURL(authAttemptStartLoginRouteName)},
		{name: "confirm login via PUT", method: http.MethodPut, target: ts.routeURL(authAttemptConfirmLoginRouteName)},
		{name: "users via POST", method: http.MethodPost, target: ts.routeURL(userDeleteRouteName)},
	}

	for _, tc := range cases {
		w := ts.sendRequest(tc.method, tc.target, nil, nil)
		ts.Equal(http.StatusMethodNotAllowed, w.Code, tc.name)
	}
}

func (ts *HandlerSuite) TestRoutes_UnknownPathIsNotFound() {
	w := ts.sendRequest(http.MethodGet, "/nope", nil, nil)
	ts.Equal(http.StatusNotFound, w.Code)
}
