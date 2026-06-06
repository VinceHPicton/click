package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (ts *HandlerSuite) TestRegisterAttemptConfirmHandler_Success() {
	// Create a register attempt first
	registerAttempt, err := ts.server.Queries.RegisterAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	// Make request to handler
	body := map[string]interface{}{
		"id":          registerAttempt.ID.String(),
		"oneTimeCode": registerAttempt.OneTimeCode,
	}
	bodyBytes, err := json.Marshal(body)
	ts.Require().NoError(err)

	confirmURL, err := ts.server.Router.Get(registerAttemptConfirmRouteName).URL()
	ts.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		confirmURL.String(),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := ts.server.registerAttemptConfirmHandler()
	handler(w, req)

	ts.Equal(http.StatusOK, w.Code)
}

func (ts *HandlerSuite) TestRegisterAttemptConfirmHandler_BadRequest() {
	// Invalid request body
	body := map[string]interface{}{
		"id": "invalid-uuid",
	}
	bodyBytes, _ := json.Marshal(body)

	confirmURL, err := ts.server.Router.Get(registerAttemptConfirmRouteName).URL()
	ts.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		confirmURL.String(),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := ts.server.registerAttemptConfirmHandler()
	handler(w, req)

	ts.Equal(http.StatusBadRequest, w.Code)
}
