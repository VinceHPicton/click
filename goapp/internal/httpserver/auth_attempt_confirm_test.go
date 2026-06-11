package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (ts *HandlerSuite) TestAuthAttemptCreateUserHandler_Success() {
	// Create an auth attempt first
	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	// Make request to handler
	body := map[string]interface{}{
		"id":          authAttempt.ID.String(),
		"oneTimeCode": authAttempt.OneTimeCode,
	}
	bodyBytes, err := json.Marshal(body)
	ts.Require().NoError(err)

	confirmURL, err := ts.server.Router.Get(authAttemptCreateUserRouteName).URL()
	ts.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		confirmURL.String(),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := ts.server.authAttemptCreateUserHandler()
	handler(w, req)

	ts.Equal(http.StatusOK, w.Code)
}

func (ts *HandlerSuite) TestAuthAttemptCreateUserHandler_BadRequest() {
	// Invalid request body
	body := map[string]interface{}{
		"id": "invalid-uuid",
	}
	bodyBytes, _ := json.Marshal(body)

	confirmURL, err := ts.server.Router.Get(authAttemptCreateUserRouteName).URL()
	ts.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		confirmURL.String(),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := ts.server.authAttemptCreateUserHandler()
	handler(w, req)

	ts.Equal(http.StatusBadRequest, w.Code)
}
