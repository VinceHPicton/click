package db

import (
	"database/sql"
	"fmt"

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

	// below results in actual IP  address array, atm its just [172.28.0.2], but obv this varies
	// addrs, err := net.LookupHost(host)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Print(addrs)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}

	return db, nil
}
