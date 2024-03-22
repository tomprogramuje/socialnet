package main

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const connStrTest = "postgres://postgres:1234@localhost:5432/test?sslmode=disable"

func TestDatabase(t *testing.T) {

	db, err := NewPostgreSQLConnection(connStrTest)
	if err != nil {
		t.Fatal(err)
	}
	clearDatabase(db)
	initializeDatabase(db)

	store := NewPostgreSQLUserStore(db)

	t.Run("creates new user Mark", func(t *testing.T) {
		username := "Mark"
		email := "test"
		password := "test"

		got, err := store.CreateUser(username, email, password)
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

		got := err.Error()
		want := fmt.Sprintf("user with id %d not found", id)

		assertError(t, got, want)
	})
	t.Run("returns Mark id", func(t *testing.T) {
		username := "Mark"

		got, err := store.GetUserByUsername(username)
		want := 1

		assertEqual(t, got.ID, want)
		assertNoError(t, err)
	})
	t.Run("stores new squeak for Mark", func(t *testing.T) {
		username := "Mark"
		squeak := "I don't believe it!"

		got, err := store.PostSqueak(username, squeak)
		want := 1

		assertEqual(t, got, want)
		assertNoError(t, err)
	})
	t.Run("get Mark's squeak", func(t *testing.T) {
		username := "Mark"

		got, err := store.GetUserSqueaks(username)
		want := []SqueakPost{{"I don't believe it!", time.Now()}}

		assertSqueaks(t, got, want)
		assertNoError(t, err)
	})
	t.Run("fetching squeaks for user with no stored squeaks", func(t *testing.T) {
		username := "Harrison"
		email := "test2"
		password := "test2"
		store.CreateUser(username, email, password)

		_, err := store.GetUserSqueaks(username)

		got := err.Error()
		want := "no squeaks found for Harrison"

		assertError(t, got, want)
	})
	t.Run("stores squeaks for Harrison and returns the userbase", func(t *testing.T) {
		username := "Harrison"
		squeak := "Great, kid, don't get cocky."
		_, err := store.PostSqueak(username, squeak)
		assertNoError(t, err)

		squeak = "Laugh it up, fuzzball!"
		_, err = store.PostSqueak(username, squeak)
		assertNoError(t, err)

		got, err := store.GetUserbase()
		want := []User{
			{1, "Mark", "test", "test", []SqueakPost{{"I don't believe it!", time.Now()}}, time.Now()},
			{2, "Harrison", "test2", "test2", []SqueakPost{{"Great, kid, don't get cocky.", time.Now()}, {"Laugh it up, fuzzball!", time.Now()}}, time.Now()},
		}

		assertUserbase(t, got, want)
		assertNoError(t, err)
	})
}
