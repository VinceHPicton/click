package sqlc_test

import (
	"time"
	"vincehpicton/click/internal/db/factory"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"

	"github.com/google/uuid"
)

func (ts *DatabaseSuite) TestRotateRefreshToken() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	_, oldToken, err := factory.FakeRefreshToken(ts.ctx, ts.queries, user.ID)
	ts.Require().NoError(err)

	newRefreshToken, err := tokens.GenerateRefreshToken()
	ts.Require().NoError(err)
	newRefreshTokenhash := tokens.HashRefreshToken(newRefreshToken)

	rotateParams := sqlc.RotateRefreshTokenParams{
		UserID:       user.ID,
		OldTokenHash: oldToken.TokenHash,
		NewTokenHash: newRefreshTokenhash,
	}

	newToken, err := ts.queries.RotateRefreshToken(ts.ctx, rotateParams)
	ts.Require().NoError(err)

	tokenRows, err := ts.queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(2, len(tokenRows))
	ts.Require().Equal(user.ID, newToken.UserID)

	_, err = ts.queries.GetValidTokenByHash(ts.ctx, oldToken.TokenHash)
	ts.Require().Error(err)

	oldTokenTBItem, err := ts.queries.GetTokenByHash(ts.ctx, oldToken.TokenHash)
	ts.NoError(err)
	ts.Equal(user.ID, oldTokenTBItem.UserID)
	ts.True(oldTokenTBItem.RevokedAt.Valid)
}

func (ts *DatabaseSuite) TestRotateRefreshToken_SucceedsIfExpiredInFuture() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	_, token, err := factory.FakeRefreshToken(ts.ctx, ts.queries, user.ID)
	ts.Require().NoError(err)

	err = ts.queries.SetRefreshTokenExpiryByID(ts.ctx, sqlc.SetRefreshTokenExpiryByIDParams{
		ExpiresAt: time.Now().Add(24 * time.Hour),
		ID:        token.ID,
	})
	ts.Require().NoError(err)

	newRefreshToken, err := tokens.GenerateRefreshToken()
	ts.Require().NoError(err)
	newRefreshTokenhash := tokens.HashRefreshToken(newRefreshToken)

	rotateParams := sqlc.RotateRefreshTokenParams{
		UserID:       user.ID,
		OldTokenHash: token.TokenHash,
		NewTokenHash: newRefreshTokenhash,
	}

	_, err = ts.queries.RotateRefreshToken(ts.ctx, rotateParams)
	ts.Require().NoError(err)

	tokenRows, err := ts.queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(2, len(tokenRows))

	_, err = ts.queries.GetValidTokenByHash(ts.ctx, token.TokenHash)
	ts.Require().Error(err)

	oldTokenTBItem, err := ts.queries.GetTokenByHash(ts.ctx, token.TokenHash)
	ts.NoError(err)
	ts.Equal(user.ID, oldTokenTBItem.UserID)
	ts.True(oldTokenTBItem.RevokedAt.Valid)
}

func (ts *DatabaseSuite) TestRotateRefreshToken_FailsWithWrongUserID() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	_, token, err := factory.FakeRefreshToken(ts.ctx, ts.queries, user.ID)
	ts.Require().NoError(err)

	invalidUserID := uuid.New()

	rotateParams := sqlc.RotateRefreshTokenParams{
		UserID:       invalidUserID,
		OldTokenHash: token.TokenHash,
		NewTokenHash: "",
	}

	_, err = ts.queries.RotateRefreshToken(ts.ctx, rotateParams)
	ts.Require().Error(err)

	tokenRows, err := ts.queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(tokenRows))

	_, err = ts.queries.GetValidTokenByHash(ts.ctx, token.TokenHash)
	ts.NoError(err)
}

func (ts *DatabaseSuite) TestRotateRefreshToken_FailsIfRevoked() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	_, token, err := factory.FakeRefreshToken(ts.ctx, ts.queries, user.ID)
	ts.Require().NoError(err)

	_, err = ts.queries.RevokeRefreshTokenByID(ts.ctx, token.ID)
	ts.Require().NoError(err)

	rotateParams := sqlc.RotateRefreshTokenParams{
		UserID:       user.ID,
		OldTokenHash: token.TokenHash,
		NewTokenHash: "",
	}

	_, err = ts.queries.RotateRefreshToken(ts.ctx, rotateParams)
	ts.Require().Error(err)

	tokenRows, err := ts.queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(tokenRows))
}

func (ts *DatabaseSuite) TestRotateRefreshToken_FailsIfExpiredInPast() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	_, token, err := factory.FakeRefreshToken(ts.ctx, ts.queries, user.ID)
	ts.Require().NoError(err)

	err = ts.queries.SetRefreshTokenExpiryByID(ts.ctx, sqlc.SetRefreshTokenExpiryByIDParams{
		ExpiresAt: time.Now().Add(-24 * time.Hour),
		ID:        token.ID,
	})
	ts.Require().NoError(err)

	rotateParams := sqlc.RotateRefreshTokenParams{
		UserID:       user.ID,
		OldTokenHash: token.TokenHash,
		NewTokenHash: "",
	}

	_, err = ts.queries.RotateRefreshToken(ts.ctx, rotateParams)
	ts.Require().Error(err)

	tokenRows, err := ts.queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(tokenRows))
}
