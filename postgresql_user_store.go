package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type PostgreSQLUserStore struct{}

func TestDB(dsName string) (*sql.DB, error) {

	db, err := sql.Open("postgres", dsName)

	if err != nil {
		log.Fatalf("couldn't connect to database, %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("couldn't verify database connection, %v", err)
	}

	return db, nil
}

func (s *PostgreSQLUserStore) PostSqueak(name, squeak string) error {
	
	return nil
}

func (s *PostgreSQLUserStore) GetUserSqueaks(name string) []string {
	return nil
}