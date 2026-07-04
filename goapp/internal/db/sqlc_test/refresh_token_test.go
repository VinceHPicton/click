package sqlc_test

import (
	"vincehpicton/click/internal/db/factory"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"
)

func (ts *DatabaseSuite) TestCreateRefreshToken() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	refreshToken := tokens.GenerateRefreshToken()
	refreshTokenHash := tokens.HashRefreshToken(refreshToken)

	p := sqlc.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
	}

	_, err = ts.queries.CreateRefreshToken(ts.ctx, p)
	ts.Require().NoError(err)

	tokenDBItem, err := ts.queries.GetValidTokenByHash(ts.ctx, refreshTokenHash)
	ts.Require().NoError(err)
	ts.Require().Equal(user.ID, tokenDBItem.UserID)

	tokenRows, err := ts.queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(tokenRows))
}
