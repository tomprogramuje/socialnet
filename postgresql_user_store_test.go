package main

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	_ "github.com/lib/pq"
)

const connStrTest = "postgres://postgres:1234@localhost:5432/test?sslmode=disable"

func TestDatabase(t *testing.T) {

	db := NewPostgreSQLConnection(connStrTest)
	clearDatabase(db)
	initializeDatabase(db)

	store := NewPostgreSQLUserStore(db)

	t.Run("creates new user Mark", func(t *testing.T) {
		name := "Mark"
		email := "test"
		password := "test"

		got, err := store.CreateUser(name, email, password)
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
		name := "Mark"

		got, err := store.GetUserByName(name)
		want := 1

		assertEqual(t, got, want)
		assertNoError(t, err)
	})
	t.Run("stores new squeak for Mark", func(t *testing.T) {
		name := "Mark"
		squeak := "I don't believe it!"

		got, err := store.PostSqueak(name, squeak)
		want := 1

		assertEqual(t, got, want)
		assertNoError(t, err)
	})
	t.Run("get Mark's squeak", func(t *testing.T) {
		user := "Mark"

		got, err := store.GetUserSqueaks(user)
		want := []string{"I don't believe it!"}

		assertSqueaks(t, got, want)
		assertNoError(t, err)
	})
	t.Run("fetching squeaks for user with no stored squeaks", func(t *testing.T) {
		username := "Harrison"
		email := "test2"
		password := "test2"
		id, _ := store.CreateUser(username, email, password)
		fmt.Println(id)

		_, err := store.GetUserSqueaks(username)

		got := err.Error()
		want := "no squeaks found for Harrison"

		assertError(t, got, want)
	})
	t.Run("stores squeaks for Harrison and returns the userbase", func(t *testing.T) {
		name := "Harrison"
		squeak := "Great, kid, don't get cocky."
		_, err := store.PostSqueak(name, squeak)
		assertNoError(t, err)

		squeak = "Laugh it up, fuzzball!"
		_, err = store.PostSqueak(name, squeak)
		assertNoError(t, err)

		got, err := store.GetUserbase()
		want := []User{
			{"Mark", "test", "test", []string{"I don't believe it!"}},
			{"Harrison", "test2", "test2", []string{"Great, kid, don't get cocky.", "Laugh it up, fuzzball!"}},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
		assertNoError(t, err)
	})
}

func assertEqual[V comparable](t testing.TB, got, want V) {
	t.Helper()
	if got != want {
		t.Error("returned value differs from expected value, got", got, "want", want)
	}
}

func assertSqueaks(t testing.TB, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("did not get correct response, got %s, want %s", got, want)
	}
}

func assertError(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q back, but wanted %q", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}
