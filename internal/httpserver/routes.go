package httpserver

import "net/http"

const (
	dbPingRouteName     = "dbping"
	dbTestRouteName     = "insertTest"
	createUserRouteName = "createUser"
)

func (s *Server) Routes() {
	s.Router.HandleFunc("/dbping", s.middlewareExample(s.handleDBping())).Name(dbPingRouteName)
	s.Router.HandleFunc("/insertTest", s.handleFakeInsert()).Methods(http.MethodGet, http.MethodPost).Name(dbTestRouteName)
	s.Router.HandleFunc("/users", s.handleCreateUser()).Methods(http.MethodPost).Name(createUserRouteName)
}
