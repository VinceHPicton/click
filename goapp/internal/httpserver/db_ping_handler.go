package httpserver

import "net/http"

func (s *Server) handleDBping() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		err := s.DB.Ping()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte("DB pinged"))
	}
}
