package db

import (
	_ "github.com/lib/pq"
	"database/sql"
	"log"
	"os"
)

// Interface to map the result of row to an interface
type RowMapper interface {
	Scan(dest ...interface{}) error
}

// Connect to database using the OS env DATABASE_URL
func Connect() *sql.DB {
	dbUrl := os.Getenv("PIKA_DATABASE_URL")
	database, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("[x] Could not open the connection to the database. Reason: %s", err.Error())
	}
	return database
}
