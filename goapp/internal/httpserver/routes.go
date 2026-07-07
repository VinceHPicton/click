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

	s.Router.HandleFunc("/auth-attempt/start/login", s.authAttemptStartLoginHandler()).Methods(http.MethodPost).Name(authAttemptStartLoginRouteName)
	s.Router.HandleFunc("/auth-attempt/start/create-user", s.authAttemptStartCreateUserHandler()).Methods(http.MethodPost).Name(authAttemptStartCreateUserRouteName)
	s.Router.HandleFunc("/auth-attempt/confirm/login", s.authAttemptConfirmLoginHandler()).Methods(http.MethodPost).Name(authAttemptConfirmLoginRouteName)
	s.Router.HandleFunc("/auth-attempt/confirm/create-user", s.authAttemptConfirmCreateUserHandler()).Methods(http.MethodPost).Name(authAttemptConfirmCreateUserRouteName)

	s.Router.HandleFunc("/auth/refresh", s.refreshHandler()).Methods(http.MethodPost).Name(refreshRouteName)

	protected := s.Router.NewRoute().Subrouter()
	protected.Use(s.authMiddleware())

	protected.HandleFunc("/users", s.userCreateHandler()).Methods(http.MethodPost).Name(userCreateRouteName)
	protected.HandleFunc("/users", s.userGetHandler()).Methods(http.MethodGet).Name(userGetRouteName)
	protected.HandleFunc("/users", s.userUpdateHandler()).Methods(http.MethodPut).Name(userUpdateRouteName)
	protected.HandleFunc("/users", s.userDeleteHandler()).Methods(http.MethodDelete).Name(userDeleteRouteName)

	protected.HandleFunc("/auth/logout", s.logoutHandler()).Methods(http.MethodPost).Name(logoutRouteName)
}
