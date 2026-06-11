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

	logoutRouteName  = "logout"
	refreshRouteName = "refresh"
)

func (s *Server) Routes() {

	s.Router.HandleFunc("/dbping", s.middlewareExample(s.handleDBping())).Name(dbPingRouteName)

	s.Router.HandleFunc("/users", s.userCreateHandler()).Methods(http.MethodPost).Name(userCreateRouteName)
	s.Router.HandleFunc("/users", s.userGetHandler()).Methods(http.MethodGet).Name(userGetRouteName)
	s.Router.HandleFunc("/users", s.userUpdateHandler()).Methods(http.MethodPut).Name(userUpdateRouteName)
	s.Router.HandleFunc("/users", s.userDeleteHandler()).Methods(http.MethodDelete).Name(userDeleteRouteName)

	s.Router.HandleFunc("/auth-attempt", s.authAttemptCreateHandler()).Methods(http.MethodPost).Name(authAttemptCreateRouteName)
	s.Router.HandleFunc("/auth-attempt/create-user", s.authAttemptCreateUserHandler()).Methods(http.MethodPost).Name(authAttemptCreateUserRouteName)

	s.Router.HandleFunc("/auth/refresh", s.refreshHandler()).Methods(http.MethodPost).Name(refreshRouteName)
	s.Router.HandleFunc("/auth/logout", s.logoutHandler()).Methods(http.MethodPost).Name(logoutRouteName)
}
