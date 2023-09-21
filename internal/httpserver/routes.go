package httpserver

const (
	dbPingRouteName = "dbping"
	dbTestRouteName = "insertTest"
)

func (s *Server) Routes() {
	s.Router.HandleFunc("/dbping", s.middlewareExample(s.handleDBping())).Name(dbPingRouteName)
	s.Router.HandleFunc("/insertTest", s.handleInsertTest()).Name(dbTestRouteName)
}
