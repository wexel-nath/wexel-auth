package database

import (
	"database/sql"
	"log"

	"github.com/wexel-nath/wexel-auth/pkg/config"

	_ "github.com/lib/pq"
)

var (
	connection *sql.DB
)

func GetConnection() *sql.DB {
	if connection == nil || connection.Ping() != nil {
		var err error
		connection, err = sql.Open("postgres", config.GetDatabaseURL())
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}
	}
	return connection
}
