package httpserver

import (
	"net/http"
)

const (
	dbPingRouteName = "dbping"
	dbTestRouteName = "insertTest"

	userCreateRouteName = "userCreate"
	userGetRouteName    = "userGet"
	userUpdateRouteName = "userUpdate"
	userDeleteRouteName = "userDelete"

	authAttemptStartLoginRouteName      = "authAttemptStartLogin"
	authAttemptStartCreateUserRouteName = "authAttemptStartCreateUser"

	authAttemptConfirmLoginRouteName      = "authAttemptConfirmLogin"
	authAttemptConfirmCreateUserRouteName = "authAttemptConfirmCreateUser"

	logoutRouteName  = "logout"
	refreshRouteName = "refresh"
)

func (s *Server) Routes() {

	s.Router.HandleFunc("/dbping", s.middlewareExample(s.handleDBping())).Name(dbPingRouteName)

	CORS := s.Router.NewRoute().Subrouter()
	CORS.Use(s.CORSMiddleware())

	CORS.Handle("/auth-attempt/start/login", s.authAttemptStartLoginHandler()).Methods(http.MethodPost).Name(authAttemptStartLoginRouteName)
	CORS.Handle("/auth-attempt/start/create-user", s.authAttemptStartCreateUserHandler()).Methods(http.MethodPost).Name(authAttemptStartCreateUserRouteName)
	CORS.Handle("/auth-attempt/confirm/login", s.authAttemptConfirmLoginHandler()).Methods(http.MethodPost).Name(authAttemptConfirmLoginRouteName)
	CORS.Handle("/auth-attempt/confirm/create-user", s.authAttemptConfirmCreateUserHandler()).Methods(http.MethodPost).Name(authAttemptConfirmCreateUserRouteName)

	CORS.HandleFunc("/auth/refresh", s.refreshHandler()).Methods(http.MethodPost).Name(refreshRouteName)

	protectedCORS := CORS.NewRoute().Subrouter()
	protectedCORS.Use(s.authMiddleware())

	protectedCORS.HandleFunc("/users", s.userCreateHandler()).Methods(http.MethodPost).Name(userCreateRouteName)
	protectedCORS.HandleFunc("/users", s.userGetHandler()).Methods(http.MethodGet).Name(userGetRouteName)
	protectedCORS.HandleFunc("/users", s.userUpdateHandler()).Methods(http.MethodPut).Name(userUpdateRouteName)
	protectedCORS.HandleFunc("/users", s.userDeleteHandler()).Methods(http.MethodDelete).Name(userDeleteRouteName)

	protectedCORS.HandleFunc("/auth/logout", s.logoutHandler()).Methods(http.MethodPost).Name(logoutRouteName)

	// A route that needs auth PLUS an extra, route-specific middleware
	protectedCORS.Handle("/users/export", Chain(s.logoutHandler(), s.authMiddleware())).Methods(http.MethodGet)
}
