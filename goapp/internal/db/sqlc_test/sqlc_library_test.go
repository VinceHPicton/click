package sqlc_test

import (
	"github.com/google/uuid"
)

func (ts *DatabaseSuite) TestExecFindsNoRowsDoesNotError() {
	err := ts.queries.BanUser(ts.ctx, uuid.New())
	ts.NoError(err)

	err = ts.queries.ConsumeAuthAttempt(ts.ctx, uuid.New())
	ts.NoError(err)
}
