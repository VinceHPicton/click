package sqlc_test

import "vincehpicton/click/internal/db/factory"

func (ts *DatabaseSuite) TestGetActiveUsersByMobile() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	users, err := ts.queries.GetActiveUsersByMobile(ts.ctx, user.Mobile)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(users))
}

func (ts *DatabaseSuite) TestGetActiveUsersByMobile_UserBanned() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	err = ts.queries.BanUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	users, err := ts.queries.GetActiveUsersByMobile(ts.ctx, user.Mobile)
	ts.Require().NoError(err)
	ts.Require().Equal(0, len(users))
}

func (ts *DatabaseSuite) TestGetActiveUsersByMobile_UserDeleted() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	err = ts.queries.SoftDeleteUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	users, err := ts.queries.GetActiveUsersByMobile(ts.ctx, user.Mobile)
	ts.Require().NoError(err)
	ts.Require().Equal(0, len(users))
}

func (ts *DatabaseSuite) TestGetActiveUsersByMobile_UserDeletedButNewExistsWithMobile() {
	var err error
	const mobile = "+447712345678"
	deletedUser, err := factory.FakeUser(ts.ctx, ts.queries, factory.WithMobile(mobile))
	ts.Require().NoError(err)

	err = ts.queries.SoftDeleteUser(ts.ctx, deletedUser.ID)
	ts.Require().NoError(err)

	undeletedUser, err := factory.FakeUser(ts.ctx, ts.queries, factory.WithMobile(mobile))
	ts.Require().NoError(err)

	users, err := ts.queries.GetActiveUsersByMobile(ts.ctx, undeletedUser.Mobile)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(users))
	ts.Require().Equal(undeletedUser.ID, users[0].ID)
}
