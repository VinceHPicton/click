package sqlc_test

import (
	"vincehpicton/click/internal/db/factory"
)

func (ts *DatabaseSuite) TestSoftDeleteUser() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	queriedUser, err := ts.queries.GetUser(ts.ctx, user.ID)
	ts.Require().NoError(err)
	ts.Require().False(queriedUser.DeletedAt.Valid)

	err = ts.queries.SoftDeleteUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	users, err := ts.queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(users))

	queriedUser, err = ts.queries.GetUser(ts.ctx, user.ID)
	ts.Require().NoError(err)
	ts.Require().True(queriedUser.DeletedAt.Valid)
}

func (ts *DatabaseSuite) TestHardDeleteUser() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	queriedUser, err := ts.queries.GetUser(ts.ctx, user.ID)
	ts.Require().NoError(err)
	ts.Require().False(queriedUser.DeletedAt.Valid)

	err = ts.queries.HardDeleteUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	users, err := ts.queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(0, len(users))
}
