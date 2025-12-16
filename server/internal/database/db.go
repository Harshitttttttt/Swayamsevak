package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect estabilishes a connection to the database provided at the dataSourceName
// and it returns an instance of the DB or an error
func Connect(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify DSN provided by the user is valid and the
	// server accessible. If the ping fails exit the program with an error.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to the Database Successfully")
	return db, nil
}
