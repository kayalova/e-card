package postgres

import (
	"database/sql"
)

//CreateConnection ...
func CreateConnection() *sql.DB {
	connStr := "user=postgres password=1 dbname=e-catalog sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	return db

}
