package httpserver

import (
	"context"
	"testing"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/testdb"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type RegisterAttemptConfirmHandlerSuite struct {
	suite.Suite
	db     testdb.TestDB
	ctx    context.Context
	server *Server
}

func TestRegisterAttemptConfirmHandlerSuite(t *testing.T) {
	suite.Run(t, new(RegisterAttemptConfirmHandlerSuite))
}

func (ts *RegisterAttemptConfirmHandlerSuite) SetupSuite() {
	ts.ctx = context.Background()

	testDB, err := testdb.Setup(ts.ctx)
	require.NoError(ts.T(), err)

	router := mux.NewRouter()

	ts.db = *testDB
	ts.server = &Server{
		DB:      ts.db.Pool,
		Queries: sqlc.New(ts.db.Pool),
		Router:  router,
	}

	ts.server.Routes()
}

func (ts *RegisterAttemptConfirmHandlerSuite) TearDownSuite() {
	ts.db.Pool.Close()
	testcontainers.CleanupContainer(ts.T(), ts.db.Container)
}

func (ts *RegisterAttemptConfirmHandlerSuite) SetupTest() {
	options := []postgres.SnapshotOption{}
	err := ts.db.Container.Snapshot(ts.ctx, options...)
	ts.Require().NoError(err)
}

func (ts *RegisterAttemptConfirmHandlerSuite) TearDownTest() {
	err := ts.db.Container.Restore(ts.ctx)
	ts.Require().NoError(err)
}
