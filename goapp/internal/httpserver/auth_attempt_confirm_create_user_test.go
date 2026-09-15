package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"vincehpicton/click/internal/db/factory"

	"github.com/google/uuid"
)

func (ts *HandlerSuite) TestAuthAttemptCreateUserHandler_Success() {
	// Create an auth attempt first
	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, "+447840195455")
	ts.Require().NoError(err)

	w := ts.callConfirmCreateUser(authAttempt.ID, authAttempt.OneTimeCode)

	ts.Equal(http.StatusOK, w.Code)

	resp := confirmCreateUserResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	ts.Require().NoError(err)

	ts.NotEqual(uuid.Nil, resp.UserID)
	ts.NotEmpty(resp.AccessToken)
	ts.NotEmpty(resp.RefreshToken)

	users, err := ts.server.Queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Len(users, 1)
	ts.Equal(users[0].ID, resp.UserID)
}

func (ts *HandlerSuite) TestAuthAttemptCreateUserHandler_BadRequest() {
	// Invalid request body
	w := ts.callConfirmCreateUserRaw(map[string]interface{}{
		"id": "invalid-uuid",
	})

	ts.Equal(http.StatusBadRequest, w.Code)
}

func (ts *HandlerSuite) TestAuthAttemptCreateUserHandler_UserAlreadyExists_AttemptStillUsed() {
	fakeUser, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, fakeUser.Mobile)
	ts.Require().NoError(err)

	w := ts.callConfirmCreateUser(authAttempt.ID, authAttempt.OneTimeCode)

	ts.Equal(http.StatusBadRequest, w.Code)

	authAttempts, err := ts.server.Queries.GetAuthAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(authAttempts))

	ts.True(authAttempts[0].UsedAt.Valid)
}

func (ts *HandlerSuite) TestAuthAttemptCreateUserHandler_UserBanned_AttemptStillUsed() {
	fakeUser, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	err = ts.server.Queries.BanUser(ts.ctx, fakeUser.ID)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, fakeUser.Mobile)
	ts.Require().NoError(err)

	w := ts.callConfirmCreateUser(authAttempt.ID, authAttempt.OneTimeCode)

	ts.Equal(http.StatusInternalServerError, w.Code)

	authAttempts, err := ts.server.Queries.GetAuthAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(authAttempts))

	ts.True(authAttempts[0].UsedAt.Valid)
}

func (ts *HandlerSuite) TestAuthAttemptCreateUserHandler_UserSoftDeleted_Success() {
	fakeUser, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, fakeUser.Mobile)
	ts.Require().NoError(err)

	ts.server.Queries.SoftDeleteUser(ts.ctx, fakeUser.ID)
	ts.Require().NoError(err)

	w := ts.callConfirmCreateUser(authAttempt.ID, authAttempt.OneTimeCode)

	ts.Require().Equal(http.StatusOK, w.Code)

	users, err := ts.server.Queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(2, len(users))

	authAttempts, err := ts.server.Queries.GetAuthAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(authAttempts))

	ts.True(authAttempts[0].UsedAt.Valid)
}

func (ts *HandlerSuite) TestAuthAttemptCreateUserHandler_ConsumedAttemptReused() {
	authAttempt, err := ts.server.Queries.AuthAttemptCreate(ts.ctx, "+447840195456")
	ts.Require().NoError(err)

	// First call consumes the auth attempt and creates the user.
	w := ts.callConfirmCreateUser(authAttempt.ID, authAttempt.OneTimeCode)
	ts.Equal(http.StatusOK, w.Code)

	// Second call with the same (now consumed) auth attempt must fail.
	reusedW := ts.callConfirmCreateUser(authAttempt.ID, authAttempt.OneTimeCode)
	ts.Require().NotEqual(http.StatusOK, reusedW.Code)

	// Only one user should have been created.
	users, err := ts.server.Queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(users))
}

func (ts *HandlerSuite) callConfirmCreateUser(id uuid.UUID, oneTimeCode int32) *httptest.ResponseRecorder {
	return ts.callConfirmCreateUserRaw(confirmCreateUserRequest{
		ID:          id,
		OneTimeCode: oneTimeCode,
	})
}

func (ts *HandlerSuite) callConfirmCreateUserRaw(body any) *httptest.ResponseRecorder {
	bodyBytes, err := json.Marshal(body)
	ts.Require().NoError(err)

	url, err := ts.server.Router.Get(authAttemptConfirmCreateUserRouteName).URL()
	ts.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		url.String(),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	ts.server.Routes()
	ts.server.Router.ServeHTTP(w, req)

	return w
}
