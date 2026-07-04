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

	authAttemptCreateRouteName     = "authAttemptCreate"
	authAttemptCreateUserRouteName = "authAttemptCreateUser"
	authAttemptLoginRouteName      = "authAttemptLogin"
	
	logoutRouteName  = "logout"
	refreshRouteName = "refresh"
)

func (s *Server) Routes() {

	s.Router.HandleFunc("/dbping", s.middlewareExample(s.handleDBping())).Name(dbPingRouteName)

	s.Router.HandleFunc("/auth-attempt", s.authAttemptCreateHandler()).Methods(http.MethodPost).Name(authAttemptCreateRouteName)
	s.Router.HandleFunc("/auth-attempt/create-user", s.authAttemptCreateUserHandler()).Methods(http.MethodPost).Name(authAttemptCreateUserRouteName)
	s.Router.HandleFunc("/auth-attempt/login", s.authAttemptLoginHandler()).Methods(http.MethodPost).Name(authAttemptLoginRouteName)

	s.Router.HandleFunc("/auth/refresh", s.refreshHandler()).Methods(http.MethodPost).Name(refreshRouteName)

	protected := s.Router.NewRoute().Subrouter()
	protected.Use(s.authMiddleware())

	protected.HandleFunc("/users", s.userCreateHandler()).Methods(http.MethodPost).Name(userCreateRouteName)
	protected.HandleFunc("/users", s.userGetHandler()).Methods(http.MethodGet).Name(userGetRouteName)
	protected.HandleFunc("/users", s.userUpdateHandler()).Methods(http.MethodPut).Name(userUpdateRouteName)
	protected.HandleFunc("/users", s.userDeleteHandler()).Methods(http.MethodDelete).Name(userDeleteRouteName)

	protected.HandleFunc("/auth/logout", s.logoutHandler()).Methods(http.MethodPost).Name(logoutRouteName)
}
