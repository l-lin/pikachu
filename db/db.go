package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
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

// Close a db without using defer
func Close(database *sql.DB) {
	err := database.Close()
	if err != nil {
		log.Printf("[x] Error when closing database. Reason: %s", err.Error())
	}
}

// Execute the given query in a transaction and return the last id of the newly created entity
func RawInsert(query string, args ...interface{}) (int, error) {
	if !strings.ContainsAny(query, "RETURNING") {
		log.Printf("[x] The query '%s' does not contains the keyword 'RETURNING' necessarily to save a new entity!", query)
	}
	database := Connect()
	tx, err := database.Begin()
	if err != nil {
		log.Printf("[x] Could not start the transaction. Reason: %s", err.Error())
		return 0, err
	}

	row := tx.QueryRow(query, args)
	var lastId int
	if err := row.Scan(&lastId); err != nil {
		log.Printf("[x] Could not fetch the last id of the newly created entity. Reason: %s", err.Error())
		if err = tx.Rollback(); err != nil {
			log.Printf("[x] Could not rollback. Reason: %s", err.Error())
		}
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
	Close(database)

	return lastId, nil
}

// Execute the update query in a transaction
func RawUpdate(query string, args ...interface{}) error {
	database := Connect()
	tx, err := database.Begin()
	if err != nil {
		log.Printf("[x] Could not start the transaction. Reason: %s", err.Error())
		return err
	}
	_, err = tx.Exec(query, args)
	if err != nil {
		log.Printf("[x] Could not update the entity. Reason: %s", err.Error())
		if err = tx.Rollback(); err != nil {
			log.Printf("[x] Could not rollback. Reason: %s", err.Error())
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
		return err
	}
	Close(database)
	return nil
}
