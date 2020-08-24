package settings

import (
	"database/sql"
	"os"
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

// GetEnvKey ...
func GetEnvKey(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); !exists {
		return value
	}

	return defaultVal
}
