package main

import (
	"database/sql"
	"log"
	"reflect"
	"slices"
	"testing"

	_ "github.com/lib/pq"
)

const connStrTest = "postgres://postgres:1234@localhost:5432/test?sslmode=disable"

func TestDatabase(t *testing.T) {

	db := NewPostgreSQLConnection(connStrTest)
	clearDatabase(db)
	initializeTestDatabase(db)

	store := NewPostgreSQLUserStore(db)

	t.Run("creates new user", func(t *testing.T) {
		name := "Mark"

		got := store.CreateUser(name)
		want := 1

		if got != want {
			t.Errorf("got wrong id back, got %d want %d", got, want)
		}
	})
	t.Run("returns user name", func(t *testing.T) {
		id := 1

		got := store.GetUserByID(id)
		want := "Mark"

		if got != want {
			t.Errorf("got wrong name back, got %s want %s", got, want)
		}

	})
	t.Run("returns not found for nonexisting user", func(t *testing.T) {
		id := 2

		got := store.GetUserByID(id)
		want := "User not found"

		if got != want {
			t.Errorf("got %s back, but wanted %s", got, want)
		}
	})
	t.Run("returns user id", func(t *testing.T) {
		name := "Mark"

		got := store.GetUserByName(name)
		want := 1

		if got != want {
			t.Errorf("got wrong id back, got %d, want %d", got, want)
		}
	})
	t.Run("stores new squeak", func(t *testing.T) {
		name := "Mark"
		squeak := "I don't believe it!"

		got := store.PostSqueak(name, squeak)

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
	})
	t.Run("fetching squeaks for user with no stored squeaks", func(t *testing.T) {
		name := "Harrison"
		store.CreateUser(name)

		got := store.GetUserSqueaks(name)
		want := []string{"No squeaks found for Harrison"}

		if !slices.Equal(got, want) {
			t.Errorf("did not get correct response, got %s, want %s", got, want)
		}
	})
	t.Run("stores squeaks for Harrison and returns the userbase", func(t *testing.T) {
		name := "Harrison"
		squeak := "Great, kid, don't get cocky."
		store.PostSqueak(name, squeak)

		squeak = "Laugh it up, fuzzball!"
		store.PostSqueak(name, squeak)

		got := store.GetUserbase()
		want := []User{
			{"Mark", []string{"I don't believe it!"}},
			{"Harrison", []string{"Great, kid, don't get cocky.", "Laugh it up, fuzzball!"}},
		}

		if !reflect.DeepEqual(got, want) { 
			t.Errorf("got %v want %v", got, want)
		}
	})
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
