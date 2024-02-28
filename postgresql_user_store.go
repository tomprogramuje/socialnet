package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type PostgreSQLUserStore struct {
	db *sql.DB
}

func TestDB(dsName string) *sql.DB {

	db, err := sql.Open("postgres", dsName)

	if err != nil {
		log.Fatalf("couldn't connect to database, %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("couldn't verify database connection, %v", err)
	}

	return db
}

func (s *PostgreSQLUserStore) CreateUser(db *sql.DB, name string) int {
	query := `INSERT INTO "user" (name)
	VALUES ($1) RETURNING id`

	var pk int
	err := db.QueryRow(query, name).Scan(&pk)
	if err != nil {
		return -1
	}

	return pk
}

func (s *PostgreSQLUserStore) GetUserByID(db *sql.DB, id int) string {
	query := `SELECT name
	FROM "user"
	WHERE id = $1
	`

	var name string
	err := db.QueryRow(query, id).Scan(&name)
	if err != nil {
		return "User not found"
	}

	return name
}

func (s *PostgreSQLUserStore) GetUserByName(db *sql.DB, name string) int {
	query := `SELECT id
	FROM "user"
	WHERE name = $1
	`

	var id int
	err := db.QueryRow(query, name).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

func (s *PostgreSQLUserStore) PostSqueak(name, squeak string) int {
	//query := `INSERT INTO squeak ()`
	return 0
}

func (s *PostgreSQLUserStore) GetUserSqueaks(name string) []string {
	return nil
}

func (s *PostgreSQLUserStore) GetUserbase() {}
