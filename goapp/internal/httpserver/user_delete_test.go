package httpserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"vincehpicton/click/internal/db/factory"
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
