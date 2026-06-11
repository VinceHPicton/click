package testdb

import (
	"context"
	"database/sql"

	"github.com/testcontainers/testcontainers-go/modules/postgres"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	postgresDBImage = "click/db"
	SQLDriver       = "pgx"
)

type TestDB struct {
	DBURL     string
	Container *postgres.PostgresContainer
	Pool      *sql.DB
}

func Setup(ctx context.Context) (*TestDB, error) {
	const dbname = "testdb"
	const dbuser = "testuser"
	const dbpassword = "testpassword"

	container, err := postgres.Run(
		ctx,
		postgresDBImage,
		postgres.WithDatabase(dbname),
		postgres.WithUsername(dbuser),
		postgres.WithPassword(dbpassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, err
	}

	// Eg: "postgres://testuser:testpassword@localhost:32769/testdb?"
	DBURL, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	pool, err := sql.Open("pgx", DBURL)
	if err != nil {
		return nil, err
	}

	testDBStruct := TestDB{
		DBURL:     DBURL,
		Container: container,
		Pool:      pool,
	}

	return &testDBStruct, nil
}
