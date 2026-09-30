package sqlc_test

import "vincehpicton/click/internal/db/factory"

func (ts *DatabaseSuite) TestGetActiveUserByMobile() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	queriedUser, err := ts.queries.GetActiveUserByMobile(ts.ctx, user.Mobile)
	ts.Require().NoError(err)
	ts.Require().Equal(user.ID, queriedUser.ID)
}

func (ts *DatabaseSuite) TestGetActiveUserByMobile_UserBanned() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	err = ts.queries.BanUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	_, err = ts.queries.GetActiveUserByMobile(ts.ctx, user.Mobile)
	ts.Error(err)
}

func (ts *DatabaseSuite) TestGetActiveUserByMobile_UserDeleted() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	err = ts.queries.SoftDeleteUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	_, err = ts.queries.GetActiveUserByMobile(ts.ctx, user.Mobile)
	ts.Error(err)
}

func (ts *DatabaseSuite) TestGetActiveUserByMobile_UserDeletedButNewExistsWithMobile() {
	var err error
	const mobile = "+447712345678"
	deletedUser, err := factory.FakeUser(ts.ctx, ts.queries, factory.WithMobile(mobile))
	ts.Require().NoError(err)

	err = ts.queries.SoftDeleteUser(ts.ctx, deletedUser.ID)
	ts.Require().NoError(err)

	undeletedUser, err := factory.FakeUser(ts.ctx, ts.queries, factory.WithMobile(mobile))
	ts.Require().NoError(err)

	queriedUser, err := ts.queries.GetActiveUserByMobile(ts.ctx, undeletedUser.Mobile)
	ts.Require().NoError(err)
	ts.Require().Equal(undeletedUser.ID, queriedUser.ID)
}

func (ts *DatabaseSuite) TestGetActiveUserByMobile_VerifyingImpossibleDBStateOfTwoActiveWithSameMobile() {
	var err error
	const mobile = "+447712345678"
	_, err = factory.FakeUser(ts.ctx, ts.queries, factory.WithMobile(mobile))
	ts.Require().NoError(err)
	_, err = factory.FakeUser(ts.ctx, ts.queries, factory.WithMobile(mobile))
	ts.Require().Error(err)
}
