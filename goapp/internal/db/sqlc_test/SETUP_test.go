package sqlc_test

import (
	"context"
	"testing"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/testdb"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type DatabaseSuite struct {
	suite.Suite
	db      testdb.TestDB
	ctx     context.Context
	queries *sqlc.Queries
}

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

func (ts *DatabaseSuite) SetupSuite() {
	ts.ctx = context.Background()

	testDB, err := testdb.Setup(ts.ctx)
	ts.Require().NoError(err)

	ts.db = *testDB
	ts.queries = sqlc.New(ts.db.Pool)
}

func (ts *DatabaseSuite) TearDownSuite() {
	ts.db.Pool.Close()
	testcontainers.CleanupContainer(ts.T(), ts.db.Container)
}

func (ts *DatabaseSuite) TestFalse() {
	ts.False(false)
}
