package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"vincehpicton/click/internal/db/factory"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"
)

func (ts *HandlerSuite) TestRefreshToken() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.server.Queries)
	ts.Require().NoError(err)

	refreshToken := tokens.GenerateRefreshToken()
	refreshTokenHash := tokens.HashRefreshToken(refreshToken)
	params := sqlc.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
	}
	err = ts.server.Queries.CreateRefreshToken(ts.ctx, params)
	ts.Require().NoError(err)

	body := map[string]interface{}{
		"refreshToken": refreshToken,
	}
	bodyBytes, err := json.Marshal(body)
	ts.Require().NoError(err)

	refreshURL, err := ts.server.Router.Get(refreshRouteName).URL()
	ts.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		refreshURL.String(),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := ts.server.refreshHandler()
	handler(w, req)

	ts.Equal(http.StatusOK, w.Code)
}
