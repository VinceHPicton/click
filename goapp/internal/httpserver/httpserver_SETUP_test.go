package httpserver

import (
	"context"
	"log"
	"os"
	"testing"
	"vincehpicton/click/internal/db/sqlc"
	"vincehpicton/click/internal/testdb"
	"vincehpicton/click/internal/tokens"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type HandlerSuite struct {
	suite.Suite
	db     testdb.TestDB
	ctx    context.Context
	server *Server
}

func TestAuthAttemptConfirmHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (ts *HandlerSuite) SetupSuite() {
	ts.ctx = context.Background()

	testDB, err := testdb.Setup(ts.ctx)
	require.NoError(ts.T(), err)

	router := mux.NewRouter()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		//TODO Address this in the future, we should be able to set env vars for the test suite, but for now we will just hardcode it
		jwtSecret = "TEST-SUPER-SECRET"
	}
	tokenManager, err := tokens.New(jwtSecret)
	if err != nil {
		log.Fatal(err)
	}

	ts.db = *testDB
	ts.server = &Server{
		DB:           ts.db.Pool,
		Queries:      sqlc.New(ts.db.Pool),
		Router:       router,
		TokenManager: tokenManager,
	}

	ts.server.Routes()
}

func (ts *HandlerSuite) TearDownSuite() {
	ts.db.Pool.Close()
	testcontainers.CleanupContainer(ts.T(), ts.db.Container)
}

func (ts *HandlerSuite) SetupTest() {
	options := []postgres.SnapshotOption{}
	err := ts.db.Container.Snapshot(ts.ctx, options...)
	ts.Require().NoError(err)
}

func (ts *HandlerSuite) TearDownTest() {
	err := ts.db.Container.Restore(ts.ctx)
	ts.Require().NoError(err)
}
