package httpserver

import "net/http"

const (
	dbPingRouteName     = "dbping"
	dbTestRouteName     = "insertTest"
	createUserRouteName = "createUser"
	updateUserRouteName = "updateUser"
)

func (s *Server) Routes() {
	s.Router.HandleFunc("/dbping", s.middlewareExample(s.handleDBping())).Name(dbPingRouteName)

	s.Router.HandleFunc("/users", s.handleCreateUser()).Methods(http.MethodPost).Name(createUserRouteName)
	s.Router.HandleFunc("/users", s.handleUpdateUser()).Methods(http.MethodPut).Name(updateUserRouteName)
}
