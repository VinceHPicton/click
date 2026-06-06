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

	registerAttemptCreateRouteName = "registerAttemptCreate"
	registerAttemptConfirmRouteName = "registerAttemptConfirm"

	loginRouteName = "login"
	logoutRouteName = "logout"
	refreshRouteName = "refresh"
)

func (s *Server) Routes() {

	s.Router.HandleFunc("/dbping", s.middlewareExample(s.handleDBping())).Name(dbPingRouteName)

	s.Router.HandleFunc("/users", s.userCreateHandler()).Methods(http.MethodPost).Name(userCreateRouteName)
	s.Router.HandleFunc("/users", s.userGetHandler()).Methods(http.MethodGet).Name(userGetRouteName)
	s.Router.HandleFunc("/users", s.userUpdateHandler()).Methods(http.MethodPut).Name(userUpdateRouteName)
	s.Router.HandleFunc("/users", s.userDeleteHandler()).Methods(http.MethodDelete).Name(userDeleteRouteName)

	s.Router.HandleFunc("/register-attempt", s.registerAttemptCreateHandler()).Methods(http.MethodPost).Name(registerAttemptCreateRouteName)
	s.Router.HandleFunc("/register-attempt/confirm", s.registerAttemptConfirmHandler()).Methods(http.MethodPost).Name(registerAttemptConfirmRouteName)

	s.Router.HandleFunc("/auth/login", s.loginHandler()).Methods(http.MethodPost).Name(loginRouteName)
	s.Router.HandleFunc("/auth/refresh", s.refreshHandler()).Methods(http.MethodPost).Name(refreshRouteName)
	s.Router.HandleFunc("/auth/logout", s.logoutHandler()).Methods(http.MethodPost).Name(logoutRouteName)
}
