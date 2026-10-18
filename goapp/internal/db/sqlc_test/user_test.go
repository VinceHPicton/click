package sqlc_test

import (
	"vincehpicton/click/internal/db/factory"
)

const (
	mobile = "+447712345678"
)

func (ts *DatabaseSuite) TestCreateUserWithMobile() {
	var err error

	newUser, err := ts.queries.CreateUserWithMobile(ts.ctx, mobile)
	ts.Require().NoError(err)

	users, err := ts.queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(1, len(users))

	queriedUser, err := ts.queries.GetUser(ts.ctx, newUser.ID)
	ts.Require().NoError(err)
	ts.Equal(newUser.Mobile, queriedUser.Mobile)
}

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

	_, err = ts.queries.GetUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	err = ts.queries.HardDeleteUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	users, err := ts.queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Equal(0, len(users))
}

func (ts *DatabaseSuite) TestBanUser() {
	var err error
	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	err = ts.queries.BanUser(ts.ctx, user.ID)
	ts.Require().NoError(err)

	user, err = ts.queries.GetUser(ts.ctx, user.ID)
	ts.Require().NoError(err)
	ts.True(user.BannedAt.Valid)
}

func (ts *DatabaseSuite) TestUnableToCreateTwoUsersWithSameMobile() {
	const mobile = "+447712345678"
	_, err := factory.FakeUser(ts.ctx, ts.queries, factory.WithMobile(mobile))
	ts.Require().NoError(err)

	_, err = factory.FakeUser(ts.ctx, ts.queries, factory.WithMobile(mobile))
	ts.Require().Error(err)
}

func (ts *DatabaseSuite) TestAbleToCreateUserWithSameMobileIfOldUserIsDeleted() {
	const mobile = "+447712345678"
	deletedUser, err := factory.FakeUser(ts.ctx, ts.queries, factory.WithMobile(mobile))
	ts.Require().NoError(err)

	err = ts.queries.SoftDeleteUser(ts.ctx, deletedUser.ID)
	ts.Require().NoError(err)

	_, err = factory.FakeUser(ts.ctx, ts.queries, factory.WithMobile(mobile))
	ts.Require().NoError(err)
}
