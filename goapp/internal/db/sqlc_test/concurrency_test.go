package sqlc_test

import (
	"sync"

	"vincehpicton/click/internal/db/factory"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/tokens"
)

// users_phone_unique_active is the only thing stopping two simultaneous
// registrations from both creating an account for one number, so the index
// doing its job under real concurrency is worth asserting directly.
func (ts *DatabaseSuite) TestCreateUserWithMobile_ConcurrentInsertsOnlyOneWins() {
	const callers = 4

	errs := make([]error, callers)
	start := make(chan struct{})
	var wg sync.WaitGroup

	for i := range errs {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, err := ts.queries.CreateUserWithMobile(ts.ctx, mobile)
			errs[i] = err
		}(i)
	}
	close(start)
	wg.Wait()

	inserted := 0
	for _, err := range errs {
		if err == nil {
			inserted++
		}
	}
	ts.Equal(1, inserted, "only one concurrent insert may claim a mobile")

	users, err := ts.queries.GetAllUsers(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(1, len(users))
}

// RotateRefreshToken revokes and inserts in one statement precisely so that a
// concurrent rotation of the same token cannot produce two live sessions. The
// losing caller's UPDATE re-checks revoked_at IS NULL once the row lock is
// released, matches nothing, and the insert is skipped.
func (ts *DatabaseSuite) TestRotateRefreshToken_ConcurrentRotationOnlyOneWins() {
	const callers = 4

	user, err := factory.FakeUser(ts.ctx, ts.queries)
	ts.Require().NoError(err)

	_, oldToken, err := factory.FakeRefreshToken(ts.ctx, ts.queries, user.ID)
	ts.Require().NoError(err)

	params := make([]sqlc.RotateRefreshTokenParams, callers)
	for i := range params {
		newToken, err := tokens.GenerateRefreshToken()
		ts.Require().NoError(err)
		params[i] = sqlc.RotateRefreshTokenParams{
			UserID:       user.ID,
			OldTokenHash: oldToken.TokenHash,
			NewTokenHash: tokens.HashRefreshToken(newToken),
		}
	}

	errs := make([]error, callers)
	start := make(chan struct{})
	var wg sync.WaitGroup

	for i := range params {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, err := ts.queries.RotateRefreshToken(ts.ctx, params[i])
			errs[i] = err
		}(i)
	}
	close(start)
	wg.Wait()

	rotated := 0
	for _, err := range errs {
		if err == nil {
			rotated++
		}
	}
	ts.Equal(1, rotated, "only one concurrent rotation may succeed")

	tokenRows, err := ts.queries.GetRefreshTokens(ts.ctx)
	ts.Require().NoError(err)
	ts.Equal(2, len(tokenRows), "the original plus exactly one replacement")

	valid := 0
	for _, row := range tokenRows {
		if !row.RevokedAt.Valid {
			valid++
		}
	}
	ts.Equal(1, valid, "exactly one live refresh token may remain")
}
