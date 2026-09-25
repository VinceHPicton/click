package httpserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"vincehpicton/click/internal/db/factory"

	"github.com/google/uuid"
)

func (ts *HandlerSuite) TestDeleteUser() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	w := ts.callDelete(user.ID.String())

	ts.Equal(http.StatusOK, w.Code)

	users, err := ts.server.Queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Len(users, 1)
	ts.Equal(true, users[0].DeletedAt.Valid)
}

func (ts *HandlerSuite) TestDeleteUser_UserDoesntExist() {
	w := ts.callDelete("non-existent-user-id")

	ts.Equal(http.StatusUnauthorized, w.Code)
}

func (ts *HandlerSuite) callDelete(userID string) *httptest.ResponseRecorder {

	token, err := ts.server.TokenManager.GenerateAccessToken(userID)
	ts.Require().NoError(err)

	url, err := ts.server.Router.Get(userDeleteRouteName).URL()
	ts.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodDelete,
		url.String(),
		bytes.NewReader([]byte("")),
	)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	ts.server.Router.ServeHTTP(w, req)

	return w
}

func (ts *HandlerSuite) TestDeleteUser_RequiresAuthentication() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	w := ts.sendRequest(http.MethodDelete, ts.routeURL(userDeleteRouteName), nil, nil)
	ts.Equal(http.StatusUnauthorized, w.Code)

	stillActive, err := ts.server.Queries.GetUser(ts.ctx, user.ID)
	ts.Require().NoError(err)
	ts.False(stillActive.DeletedAt.Valid)
}

// SoftDeleteUser is an :exec, so a well-formed UUID that matches no row updates
// nothing and reports success. The endpoint cannot be used to probe for which
// accounts exist, which is the behaviour worth having.
func (ts *HandlerSuite) TestDeleteUser_UnknownButValidUUIDSucceedsSilently() {
	existing, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	w := ts.callDelete(uuid.New().String())
	ts.Equal(http.StatusOK, w.Code)

	untouched, err := ts.server.Queries.GetUser(ts.ctx, existing.ID)
	ts.Require().NoError(err)
	ts.False(untouched.DeletedAt.Valid)
}

func (ts *HandlerSuite) TestDeleteUser_SecondDeleteStillSucceeds() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	first := ts.callDelete(user.ID.String())
	ts.Require().Equal(http.StatusOK, first.Code)

	second := ts.callDelete(user.ID.String())
	ts.Equal(http.StatusOK, second.Code)
}

// Deleting an account leaves its refresh token rows in place and un-revoked.
// Nothing can rotate them afterwards, because RefreshSession rejects a deleted
// user, but the rows survive. Pinned so a future change to revoke-on-delete is
// a deliberate decision rather than an accident.
func (ts *HandlerSuite) TestDeleteUser_LeavesRefreshTokensUnrevoked() {
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken, _, err := factory.FakeRefreshToken(ts.ctx, ts.server.Queries, user.ID)
	ts.Require().NoError(err)

	deleted := ts.callDelete(user.ID.String())
	ts.Require().Equal(http.StatusOK, deleted.Code)

	tokenRows, err := ts.server.Queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Len(tokenRows, 1)
	ts.False(tokenRows[0].RevokedAt.Valid)

	ts.Equal(http.StatusBadRequest, ts.callRefresh(refreshToken).Code, "but it can no longer be rotated")
}
