package httpserver

import (
	"net/http"
)

func (s *Server) handledbping() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		err := s.DB.Ping()
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte("DB pinged"))
	}
}
