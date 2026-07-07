package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
)

func (ts *HandlerSuite) TestStartLogin() {
	const mobile = "+447840195452"
	body, err := json.Marshal(authAttemptStartRequest{
		Mobile: mobile,
	})
	ts.Require().NoError(err)

	url, err := ts.server.Router.Get(authAttemptStartLoginRouteName).URL()
	ts.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		url.String(),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler := ts.server.authAttemptStartLoginHandler()
	handler(w, req)

	ts.Require().Equal(http.StatusOK, w.Code)

	resp := authAttemptStartResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	ts.Require().NoError(err)

	ts.Require().NotEmpty(resp.ID)
	_, err = uuid.Parse(resp.ID.String())
	ts.Require().NoError(err)
	ts.Require().Equal(6, len(resp.OneTimeCode))

	attempts, err := ts.server.Queries.GetAuthAttempts(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(len(attempts), 1)
	ts.Require().Equal(attempts[0].Mobile, mobile)
}
