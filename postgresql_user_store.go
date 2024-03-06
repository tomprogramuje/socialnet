package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const connStrProd = "postgres://postgres:1234@localhost:5432/postgres?sslmode=disable"

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

func (s *PostgreSQLUserStore) CreateUser(name string) (int, error) {
	query := `INSERT INTO "user" (name) VALUES ($1) RETURNING id`

	var id int
	err := s.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %w", err)
	}

	return id, nil
}

func (s *PostgreSQLUserStore) GetUserByID(id int) (string, error) {
	query := `SELECT name FROM "user" WHERE id = $1`

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

func (s *PostgreSQLUserStore) GetUserByName(name string) (int, error) {
	query := `SELECT id FROM "user"	WHERE name = $1`

	var id int
	err := s.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("no user of that name (%s) found", name)
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

func (s *PostgreSQLUserStore) GetUserSqueaks(name string) ([]string, error) {
	user_id, err := s.GetUserByName(name)
	if err != nil {
		return nil, fmt.Errorf("no user of that name (%s) found", name)
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
		return nil, fmt.Errorf("no squeaks found for %s", name)
	}

	return squeaks, nil
}

func (s *PostgreSQLUserStore) GetUserbase() ([]User, error) {
	query := `SELECT name, text FROM "user" u JOIN "squeak" s 
		ON u.id = s.user_id ORDER BY u.id, s.id;`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetUserbase: %w", err)
	}

	var userbase []User

	for rows.Next() {
		var name, squeak string
		if err := rows.Scan(&name, &squeak); err != nil {
			return nil, fmt.Errorf("GetUserbase: %w", err)
		}

		userExists := false
		for i := range userbase {
			if userbase[i].Name == name {
				userbase[i].Squeaks = append(userbase[i].Squeaks, squeak)
				userExists = true
				break
			}
		}

		if !userExists {
			userbase = append(userbase, User{Name: name, Squeaks: []string{squeak}})
		}
	}

	return userbase, nil
}
