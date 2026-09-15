package httpserver

import (
	"encoding/json"
	"net/http"
)

const maxBodyBytes = 1 << 20

// Validator lets a request type declare its own field-level checks. decodeJSON
// runs Valid after a successful decode when the type implements this.
type Validator interface {
	Valid() bool
}

// decodeJSON reads a single JSON object from the request body into T.
//
// It bounds the body, rejects unknown fields and trailing content, and runs
// Valid when T is a Validator. On any failure it writes the HTTP error itself
// and returns ok=false, so callers only need to return.
func decodeJSON[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var req T

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "malformed request body", http.StatusBadRequest)
		return req, false
	}

	if decoder.More() {
		http.Error(w, "malformed request body", http.StatusBadRequest)
		return req, false
	}

	if v, ok := any(req).(Validator); ok && !v.Valid() {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return req, false
	}

	return req, true
}
