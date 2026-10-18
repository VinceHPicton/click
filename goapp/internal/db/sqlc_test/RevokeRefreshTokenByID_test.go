package sqlc_test

import (
	"database/sql"
	"errors"

	"vincehpicton/click/internal/db/factory"

	"github.com/google/uuid"
)

func (ts *DatabaseSuite) TestRevokeRefreshTokenByID() {
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	_, token, err := factory.FakeRefreshToken(ts.ctx, ts.queries, user.ID)
	ts.Require().NoError(err)

	revoked, err := ts.queries.RevokeRefreshTokenByID(ts.ctx, token.ID)
	ts.Require().NoError(err)
	ts.Require().True(revoked.RevokedAt.Valid)
	ts.Equal(token.ID, revoked.ID)

	_, err = ts.queries.GetValidTokenByHash(ts.ctx, token.TokenHash)
	ts.Require().Error(err)
	ts.True(errors.Is(err, sql.ErrNoRows))
}

func (ts *DatabaseSuite) TestRevokeRefreshTokenByID_UnknownIDReturnsNoRows() {
	_, err := ts.queries.RevokeRefreshTokenByID(ts.ctx, uuid.New())
	ts.Require().Error(err)
	ts.True(errors.Is(err, sql.ErrNoRows))
}

// The query has no revoked_at IS NULL guard, so a second revoke still matches
// the row and re-stamps it. service.Logout relies on GetValidTokenByHash to
// reject the replay before reaching here.
func (ts *DatabaseSuite) TestRevokeRefreshTokenByID_HasNoAlreadyRevokedGuard() {
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	_, token, err := factory.FakeRefreshToken(ts.ctx, ts.queries, user.ID)
	ts.Require().NoError(err)

	_, err = ts.queries.RevokeRefreshTokenByID(ts.ctx, token.ID)
	ts.Require().NoError(err)

	reRevoked, err := ts.queries.RevokeRefreshTokenByID(ts.ctx, token.ID)
	ts.Require().NoError(err)
	ts.True(reRevoked.RevokedAt.Valid)
}
