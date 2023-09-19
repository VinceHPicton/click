package httpserver

func (s *Server) routes() {

	s.Router.HandleFunc("/dbping", s.middlewareExample(s.handledbping()))
}
