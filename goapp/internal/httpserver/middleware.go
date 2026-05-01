package httpserver

import "net/http"

func (s *Server) middlewareExample(prevHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		prevHandler(w, r)
	}
}

func (s *Server) corsMiddleware(prevHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		prevHandler(w, r)
	}
}
