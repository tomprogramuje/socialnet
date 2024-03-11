package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const connStrProd = "postgres://postgres:1234@localhost:5432/postgres?sslmode=disable"

func initializeDatabase(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS "user" (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) NOT NULL UNIQUE,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS "squeak" (
		id SERIAL PRIMARY KEY,
		user_id INT,
		text VARCHAR(255),
		FOREIGN KEY (user_id) REFERENCES "user"(id)
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("error initializing database: %v", err)
	}
}

func clearDatabase(db *sql.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS squeak; DROP TABLE IF EXISTS "user";`)
	if err != nil {
		log.Fatalf("error dropping table: %v", err)
	}
}

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

func (s *PostgreSQLUserStore) CreateUser(username, email, password string) (int, error) {
	query := `INSERT INTO "user" (username, email, password) VALUES ($1, $2, $3) RETURNING id`

	var id int
	err := s.db.QueryRow(query, username, email, password).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %w", err)
	}

	return id, nil
}

func (s *PostgreSQLUserStore) GetUserByID(id int) (string, error) {
	query := `SELECT username FROM "user" WHERE id = $1`

	var name string
	err := s.db.QueryRow(query, id).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user with id %d not found", id)
		}
		return "", fmt.Errorf("GetUserByID: %w", err)
	}

	return name, nil
}

func (s *PostgreSQLUserStore) GetUserByName(username string) (int, error) {
	query := `SELECT id FROM "user"	WHERE username = $1`

	var id int
	err := s.db.QueryRow(query, username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("no user with that username (%s) found", username)
		}
		return 0, fmt.Errorf("GetUserByName: %w", err)
	}

	return id, nil
}

func (s *PostgreSQLUserStore) PostSqueak(name, squeak string) (int, error) {
	user_id, err := s.GetUserByName(name)
	if err != nil {
		return 0, fmt.Errorf("error trying to post new squeak: %s", err)
	}
	query := `INSERT INTO squeak (user_id, text) VALUES ($1, $2) RETURNING id`

	var id int
	err = s.db.QueryRow(query, user_id, squeak).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("PostSqueak: %w", err)
	}

	return id, nil
}

func (s *PostgreSQLUserStore) GetUserSqueaks(username string) ([]string, error) {
	user_id, err := s.GetUserByName(username)
	if err != nil {
		return nil, fmt.Errorf("no user with that username (%s) found", username)
	}
	query := `SELECT text FROM squeak WHERE user_id = $1`

	var squeaks []string
	rows, err := s.db.Query(query, user_id)
	if err != nil {
		return nil, fmt.Errorf("GetUserSqueaks: %w", err)
	}

	defer rows.Close()

	var squeak string
	for rows.Next() {
		if err := rows.Scan(&squeak); err != nil {
			return nil, fmt.Errorf("GetUserSqueaks: %w", err)
		}
		squeaks = append(squeaks, squeak)
	}

	if len(squeaks) == 0 {
		return nil, fmt.Errorf("no squeaks found for %s", username)
	}

	return squeaks, nil
}

func (s *PostgreSQLUserStore) GetUserbase() ([]User, error) {
	query := `SELECT username, email, password, text FROM "user" u JOIN "squeak" s 
		ON u.id = s.user_id ORDER BY u.id, s.id;`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetUserbase: %w", err)
	}

	var userbase []User

	for rows.Next() {
		var username, email, password, squeak string
		if err := rows.Scan(&username, &email, &password, &squeak); err != nil {
			return nil, fmt.Errorf("GetUserbase: %w", err)
		}

		userExists := false
		for i := range userbase {
			if userbase[i].Username == username {
				userbase[i].Squeaks = append(userbase[i].Squeaks, squeak)
				userExists = true
				break
			}
		}

		if !userExists {
			userbase = append(userbase, User{Username: username, Email: email, Password: password, Squeaks: []string{squeak}})
		}
	}

	return userbase, nil
}
