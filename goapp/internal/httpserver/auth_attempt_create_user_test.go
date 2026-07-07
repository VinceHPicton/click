package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"vincehpicton/click/internal/db/factory"
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

	confirmURL, err := ts.server.Router.Get(authAttemptConfirmCreateUserRouteName).URL()
	ts.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		confirmURL.String(),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := ts.server.authAttemptConfirmCreateUserHandler()
	handler(w, req)

	ts.Equal(http.StatusOK, w.Code)
}

func (ts *HandlerSuite) TestAuthAttemptCreateUserHandler_BadRequest() {
	// Invalid request body
	body := map[string]interface{}{
		"id": "invalid-uuid",
	}
	bodyBytes, _ := json.Marshal(body)

	confirmURL, err := ts.server.Router.Get(authAttemptConfirmCreateUserRouteName).URL()
	ts.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		confirmURL.String(),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := ts.server.authAttemptConfirmCreateUserHandler()
	handler(w, req)

	ts.Equal(http.StatusBadRequest, w.Code)
}

func (ts *HandlerSuite) TestAuthAttemptCreateUserHandler_UserAlreadyExists() {
	const phoneNumber = "+447840195455"

	fakeUser, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, fakeUser.Mobile)
	ts.Require().NoError(err)

	// Make request to handler
	body := map[string]interface{}{
		"id":          authAttempt.ID.String(),
		"oneTimeCode": authAttempt.OneTimeCode,
	}
	bodyBytes, err := json.Marshal(body)
	ts.Require().NoError(err)

	confirmURL, err := ts.server.Router.Get(authAttemptConfirmCreateUserRouteName).URL()
	ts.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		confirmURL.String(),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := ts.server.authAttemptConfirmCreateUserHandler()
	handler(w, req)

	ts.Equal(http.StatusInternalServerError, w.Code)

	// TODO: if user already exists, auth attempt isn't consumed.
	// authAttempts, err := ts.server.Queries.GetAuthAttempts(ts.ctx)
	// ts.Require().NoError(err)
	// ts.Require().Equal(1, len(authAttempts))

	// ts.True(authAttempts[0].UsedAt.Valid)
}
