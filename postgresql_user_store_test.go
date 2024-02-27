package main

import (
	"testing"

	_ "github.com/lib/pq"
)

const connStrTest = "postgres://postgres:1234@localhost:5432/test?sslmode=disable"

func TestDatabase(t *testing.T) {
	
	db, err := TestDB(connStrTest)
	if err != nil {
		t.Errorf("error connecting to database: %v", err)
	}
	query := `CREATE TABLE IF NOT EXISTS test_user (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL
	)`

	_, err = db.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
	
	/*t.Run("stores new squeak", func(t *testing.T) {
		user := "Mark"
		squeak := "I don't believe it!"
		store := PostgreSQLUserStore{}

		err := store.PostSqueak(user, squeak)
		if err != nil {
			t.Error("error inserting data into database", err)
		}
	})
	t.Run("get user squeak", func(t *testing.T) {
		user := "Mark"
		store := PostgreSQLUserStore{}

		got := store.GetUserSqueaks(user)
		want := []string{"I don't believe it!"}

		if !slices.Equal(got, want) {
			t.Errorf("did not get correct reponse, got %s, want %s", got, want)
		}
	})*/

	_, err = db.Exec("DROP TABLE IF EXISTS test_user")
	if err != nil {
		t.Fatalf("error dropping table: %v", err)
	}
}
