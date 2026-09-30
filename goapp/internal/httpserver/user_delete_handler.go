package httpserver

import (
	"net/http"
)

func (s *Server) userDeleteHandler() http.HandlerFunc {
	// TODO: leaving these here to just illustrate the pattern of structs scoped to the handler
	type request struct {
	}
	type response struct {
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		auth, ok := AuthFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusInternalServerError)
			return
		}

		err = s.Queries.SoftDeleteUser(r.Context(), auth.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
