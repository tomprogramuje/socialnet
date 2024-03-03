package main

import (
	"database/sql"
	"fmt"
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

	t.Run("creates new user Mark", func(t *testing.T) {
		name := "Mark"
		
		got, err := store.CreateUser(name)
		want := 1

		assertEqual(t, got, want)
		assertNoError(t, err)
	})
	t.Run("returns user id 1 name", func(t *testing.T) {
		id := 1
		
		got, err := store.GetUserByID(id)
		want := "Mark"

		assertEqual(t, got, want)
		assertNoError(t, err)
	})
	t.Run("returns not found for nonexisting user", func(t *testing.T) {
		id := 2
		_, err := store.GetUserByID(id)
		want := fmt.Errorf("User with id %d not found", id)

		if err.Error() != want.Error() {
			t.Errorf("got %q back, but wanted %q", err, want)
		}
	})
	t.Run("returns Mark id", func(t *testing.T) {
		name := "Mark"
		
		got, err := store.GetUserByName(name)
		want := 1

		assertEqual(t, got, want)
		assertNoError(t, err)
	})
	t.Run("stores new squeak for Mark", func(t *testing.T) {
		name := "Mark"
		squeak := "I don't believe it!"

		got, _ := store.PostSqueak(name, squeak)
		want := 1

		assertEqual(t, got, want)
	})
	t.Run("get Mark's squeak", func(t *testing.T) {
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
		_, err := store.PostSqueak(name, squeak)
		assertNoError(t, err)

		squeak = "Laugh it up, fuzzball!"
		_, err = store.PostSqueak(name, squeak)
		assertNoError(t, err)

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

func assertEqual[V comparable](t testing.TB, got, want V) {
	t.Helper()
	if got != want {
		t.Error("returned values differ, got", got, "want", want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
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
