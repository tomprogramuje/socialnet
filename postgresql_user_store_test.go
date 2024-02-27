package main

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

const connStrTest = "postgres://postgres:1234@localhost:5432/test?sslmode=disable"

func TestDatabase(t *testing.T) {

	db, err := TestDB(connStrTest)
	if err != nil {
		t.Errorf("error connecting to database: %v", err)
	}

	clearDatabase(db)

	initializeTestDatabase(db)

	store := PostgreSQLUserStore{db: db}

	t.Run("creates new user", func(t *testing.T) {
		name := "Mark"

		got := store.CreateUser(db, name)
		want := 1

		if got != want {
			t.Errorf("got wrong id back, got %d want %d", got, want)
		}
	})
	t.Run("returns user name", func(t *testing.T) {
		id := 1

		got := store.GetUserByID(db, id)
		want := "Mark"

		if got != want {
			t.Errorf("got wrong name back, got %s want %s", got, want)
		}

	})
	/*t.Run("stores new squeak", func(t *testing.T) {
		user := "Mark"
		squeak := "I don't believe it!"

		got, err := store.PostSqueak(user, squeak)
		if err != nil {
			t.Error("error inserting data into database", err)
		}

		want := 1

		if got != want {
			t.Errorf("got wrong id back, got %d want %d", got, want)
		}
	})
	t.Run("get user squeak", func(t *testing.T) {
		user := "Mark"

		got := store.GetUserSqueaks(user)
		want := []string{"I don't believe it!"}

		if !slices.Equal(got, want) {
			t.Errorf("did not get correct reponse, got %s, want %s", got, want)
		}
	})*/
}

func clearDatabase(db *sql.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS squeak; DROP TABLE IF EXISTS "user";`)
	if err != nil {
		log.Fatalf("error dropping table: %v", err)
	}
}

func initializeTestDatabase(db *sql.DB) {
	query := `CREATE TABLE "user" (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL
	);
	CREATE TABLE "squeak" (
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