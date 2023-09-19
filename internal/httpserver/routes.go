package httpserver

func (s *Server) Routes() {

	s.Router.HandleFunc("/dbping", s.middlewareExample(s.handleDBping()))
	s.Router.HandleFunc("/insertTest", s.handleInsertTest())
}
