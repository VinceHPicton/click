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
)

func (s *Server) Routes() {
	s.Router.HandleFunc("/dbping", s.middlewareExample(s.handleDBping())).
		Name(dbPingRouteName)

	s.Router.HandleFunc("/users", s.corsMiddleware(s.userCreateHandler())).
		Methods(http.MethodPost).Name(userCreateRouteName)

	s.Router.HandleFunc("/users", s.corsMiddleware(s.userGetHandler())).
		Methods(http.MethodGet).Name(userGetRouteName)

	s.Router.HandleFunc("/users", s.corsMiddleware(s.userUpdateHandler())).
		Methods(http.MethodPut).Name(userUpdateRouteName)

	s.Router.HandleFunc("/users", s.corsMiddleware(s.userDeleteHandler())).
		Methods(http.MethodDelete).Name(userDeleteRouteName)

	s.Router.HandleFunc("/register-attempt", s.corsMiddleware(s.registerAttemptCreateHandler())).
		Methods(http.MethodPost).Name(registerAttemptCreateRouteName)
}
