package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const connStrProd = "postgres://postgres:1234@localhost:5432/postgres?sslmode=disable"

func initializeDatabase(db *sql.DB) { 
	query := `CREATE TABLE IF NOT EXISTS "user" (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(60) NOT NULL,
		createdAt TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS "squeak" (
		id SERIAL PRIMARY KEY,
		user_id INT,
		text VARCHAR(255),
		createdAt TIMESTAMP,
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

func NewPostgreSQLConnection(dsName string) (*sql.DB, error) {

	db, err := sql.Open("postgres", dsName)

	if err != nil {
		return nil, fmt.Errorf("couldn't connect to database, %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("couldn't verify database connection, %v", err)
	}

	return db, nil
}

func (s *PostgreSQLUserStore) CreateUser(username, email, password string) (int, error) {
	query := `INSERT INTO "user" (username, email, password, createdAt) VALUES ($1, $2, $3, $4) RETURNING id`

	var id int
	err := s.db.QueryRow(query, username, email, password, time.Now()).Scan(&id)
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

func (s *PostgreSQLUserStore) GetUserByUsername(username string) (*User, error) {
	query := `SELECT * FROM "user"	WHERE username = $1`

	user := new(User)
	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no user with that username (%s) found", username)
		}
		return nil, fmt.Errorf("GetUserByName: %w", err)
	}

	return user, nil
}

func (s *PostgreSQLUserStore) PostSqueak(name, squeak string) (int, error) {
	user, err := s.GetUserByUsername(name)
	if err != nil {
		return 0, fmt.Errorf("error trying to post new squeak: %s", err)
	}
	query := `INSERT INTO squeak (user_id, text, createdAt) VALUES ($1, $2, $3) RETURNING id`

	var id int
	err = s.db.QueryRow(query, user.ID, squeak, time.Now()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("PostSqueak: %w", err)
	}

	return id, nil
}

func (s *PostgreSQLUserStore) GetUserSqueaks(username string) ([]SqueakPost, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("no user with that username (%s) found", username)
	}
	query := `SELECT text, createdAt FROM squeak WHERE user_id = $1`

	var squeaks []SqueakPost
	rows, err := s.db.Query(query, user.ID)
	if err != nil {
		return nil, fmt.Errorf("GetUserSqueaks: %w", err)
	}

	defer rows.Close()

	var squeak SqueakPost
	for rows.Next() {
		if err := rows.Scan(&squeak.Text, &squeak.CreatedAt); err != nil {
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
	query := `SELECT u.id, username, email, password, s.text, s.createdAt, u.createdAt FROM "user" u JOIN "squeak" s 
		ON u.id = s.user_id ORDER BY u.id, s.id;`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetUserbase: %w", err)
	}

	var userbase []User

	for rows.Next() {
		user := new(User)
		var squeak SqueakPost
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &squeak.Text, &squeak.CreatedAt, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetUserbase: %w", err)
		}

		userExists := false
		for i := range userbase {
			if userbase[i].Username == user.Username {
				userbase[i].Squeaks = append(userbase[i].Squeaks, squeak)
				userExists = true
				break
			}
		}

		if !userExists {
			userbase = append(userbase, User{ID: user.ID, Username: user.Username, Email: user.Email, Password: user.Password, Squeaks: []SqueakPost{squeak}})
		}
	}

	return userbase, nil
}
