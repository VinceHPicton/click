package httpserver

import (
	"encoding/json"
	"net/http"
	"vincehpicton/click/internal/db/sqlc"

	"github.com/google/uuid"
)

func (s *Server) registerAttemptConfirmHandler() http.HandlerFunc {

	type request struct {
		ID uuid.UUID `json:"id"`
		OneTimeCode int32 `json:"one_time_code"`
	}

	type response struct {
	}

	return func(w http.ResponseWriter, r *http.Request) {

		queries := sqlc.New(s.DB)

		confirmRegisterAttemptParams := sqlc.RegisterAttemptConfirmParams{}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&confirmRegisterAttemptParams)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		registerAttemptId, err := queries.RegisterAttemptConfirm(r.Context(), confirmRegisterAttemptParams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if registerAttemptId == uuid.Nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
	}
}
