package httpserver

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"vincehpicton/click/internal/service"
)

// writeJSON encodes body as the response, with the given status.
//
// The header has to be written before the body, and once bytes are on the wire
// the status can't be changed, so an encoding failure part-way through can only
// be logged rather than turned into a 500.
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("httpserver: encode response body error: %v", err)
	}
}

// writeError is the single place where a service error becomes an HTTP status.
//
// Keeping the mapping here rather than in each handler is what stops the same
// failure from being a 400 on one route and a 500 on another.
//
// Unrecognised errors are logged in full and reported to the client as a bare
// 500: handlers used to pass err.Error() straight through, which leaked raw
// Postgres messages to callers.
func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidAttempt), errors.Is(err, service.ErrInvalidCode):
		// Same response for both, so a caller can't tell a bad attempt ID from
		// a bad code.
		http.Error(w, "invalid or expired one-time code", http.StatusUnauthorized)

	case errors.Is(err, service.ErrUserNotFound):
		// Intentionally vague: banned and soft-deleted users land here too, and
		// which of the three it is isn't the client's business.
		http.Error(w, "account not found", http.StatusNotFound)

	case errors.Is(err, service.ErrUserExists):
		http.Error(w, "account already exists", http.StatusBadRequest)

	default:
		log.Printf("httpserver: unhandled error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
