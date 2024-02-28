package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgreSQLUserStore struct {
	db *sql.DB
}

func NewPostgreSQLUserStore(db *sql.DB) *PostgreSQLUserStore {
	return &PostgreSQLUserStore{db: db}
}

func NewPostgreSQLConnection(dsName string) *sql.DB {

	db, err := sql.Open("postgres", dsName)

	if err != nil {
		log.Fatalf("couldn't connect to database, %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("couldn't verify database connection, %v", err)
	}

	return db
}

func (s *PostgreSQLUserStore) CreateUser(name string) int {
	query := `INSERT INTO "user" (name)
	VALUES ($1) RETURNING id`

	var pk int
	err := s.db.QueryRow(query, name).Scan(&pk)
	if err != nil {
		return -1
	}

	return pk
}

func (s *PostgreSQLUserStore) GetUserByID(id int) string {
	query := `SELECT name
	FROM "user"
	WHERE id = $1
	`

	var name string
	err := s.db.QueryRow(query, id).Scan(&name)
	if err != nil {
		return "User not found"
	}

	return name
}

func (s *PostgreSQLUserStore) GetUserByName(name string) int {
	query := `SELECT id
	FROM "user"
	WHERE name = $1
	`

	var id int
	err := s.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1
		}
		log.Fatal(err)
	}

	return id
}

func (s *PostgreSQLUserStore) PostSqueak(name, squeak string) int {
	user_id := s.GetUserByName(name)
	query := `INSERT INTO squeak (user_id, text)
	VALUES ($1, $2) RETURNING id`

	var id int
	err := s.db.QueryRow(query, user_id, squeak).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

func (s *PostgreSQLUserStore) GetUserSqueaks(name string) []string {
	user_id := s.GetUserByName(name)
	query := `SELECT text FROM squeak WHERE user_id = $1`

	var squeaks []string
	rows, err := s.db.Query(query, user_id)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var squeak string
	for rows.Next() {
		err := rows.Scan(&squeak)
		if err != nil {
			log.Fatal(err)
		}
		squeaks = append(squeaks, squeak)
	}

	if len(squeaks) == 0 {
		return []string{fmt.Sprintf("No squeaks found for %s", name)}
	}

	return squeaks
}

func (s *PostgreSQLUserStore) GetUserbase() []User {
	// iterating through all the users in db
	// then iterate through all the squeaks of the specific user 
	// pass the data to User struct 
	// return a slice of all User structs
	//s.GetUserByID(id int)
	
	
	return []User{{"Mark", []string{"I don't believe it!"}}}
}
