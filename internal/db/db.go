package db

import (
	"database/sql"
	"fmt"
	"net"

	_ "github.com/lib/pq"
)

const (
	host     = "dbserver"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func Open() (*sql.DB, error) {

	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		addrs[0], port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}

	return db, nil
}
