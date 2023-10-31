package httpserver

import "net/http"

const (
	dbPingRouteName     = "dbping"
	dbTestRouteName     = "insertTest"
	createUserRouteName = "createUser"
	getUserRouteName    = "getUser"
	updateUserRouteName = "updateUser"
	deleteUserRouteName = "deleteUser"
)

func (s *Server) Routes() {
	s.Router.HandleFunc("/dbping", s.middlewareExample(s.handleDBping())).Name(dbPingRouteName)

	s.Router.HandleFunc("/users", s.handleCreateUser()).Methods(http.MethodPost).Name(createUserRouteName)
	s.Router.HandleFunc("/users", s.handleGetUser()).Methods(http.MethodGet).Name(getUserRouteName)
	s.Router.HandleFunc("/users", s.handleUpdateUser()).Methods(http.MethodPut).Name(updateUserRouteName)
	s.Router.HandleFunc("/users", s.handleDeleteUser()).Methods(http.MethodDelete).Name(deleteUserRouteName)
}
