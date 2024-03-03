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

func (s *PostgreSQLUserStore) CreateUser(name string) (int, error) {
	query := `INSERT INTO "user" (name) VALUES ($1) RETURNING id`

	var id int
	err := s.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		return -1, nil // todo
	}

	return id, nil
}

func (s *PostgreSQLUserStore) GetUserByID(id int) (string, error) {
	query := `SELECT name FROM "user" WHERE id = $1`

	var name string
	err := s.db.QueryRow(query, id).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("User with id %d not found", id)
		}
		log.Fatal(err)
	}

	return name, nil
}

func (s *PostgreSQLUserStore) GetUserByName(name string) (int, error) {
	query := `SELECT id FROM "user"	WHERE name = $1`

	var id int
	err := s.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, fmt.Errorf("no user of that name (%s) found", name)
		}
		log.Fatal(err)
	}

	return id, nil
}

func (s *PostgreSQLUserStore) PostSqueak(name, squeak string) (int, error) {
	user_id, err := s.GetUserByName(name)
	if err != nil {
		return 0, err // todo
	}
	query := `INSERT INTO squeak (user_id, text) VALUES ($1, $2) RETURNING id`

	var id int
	err = s.db.QueryRow(query, user_id, squeak).Scan(&id)
	if err != nil {
		return -1, nil // todo
	}

	return id, nil
}

func (s *PostgreSQLUserStore) GetUserSqueaks(name string) []string {
	user_id, err := s.GetUserByName(name)
	if err != nil {
		return []string{fmt.Sprintf("No user of that name (%s) found", name)}
	}
	query := `SELECT text FROM squeak WHERE user_id = $1`

	var squeaks []string
	rows, err := s.db.Query(query, user_id)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var squeak string
	for rows.Next() {
		if err := rows.Scan(&squeak); err != nil {
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
	query := `SELECT name, text FROM "user" u JOIN "squeak" s 
		ON u.id = s.user_id ORDER BY u.id, s.id;`

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

	return userbase
}

/*func handleErrors[T any](err error, returnValue T) returnValue {
	if err != nil {
		if err == sql.ErrNoRows {
			return errMsg
		}
		log.Fatal(err)
	}

	return ""
}*/

// todo: errorHelper func - probably generic func returning different types of data, db fetching helpers?
