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

	var id int
	err := s.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

func (s *PostgreSQLUserStore) GetUserByID(id int) string {
	query := `SELECT name
	FROM "user"
	WHERE id = $1
	`

	var name string
	err := s.db.QueryRow(query, id).Scan(&name)
	if err != nil {
		return "User not found" // shouldnť check for sql.ErrNoRows?
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
	if err != nil { // shouldnť check for sql.ErrNoRows?
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
	query := `SELECT name, text 
	FROM "user" u
	JOIN "squeak" s 
	ON u.id = s.user_id
	ORDER BY u.id;`

	rows, err := s.db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}

	var userbase []User
	for rows.Next() {
		var name, squeak string
		if err := rows.Scan(&name, &squeak); err != nil {
			log.Fatal(err)
		}
		user := &User{Name: name, Squeaks: []string{squeak}}
		userbase = append(userbase, *user)
	}

	return userbase
}
