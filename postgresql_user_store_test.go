package main

import (
	"slices"
	"testing"

	_ "github.com/lib/pq"
)

const connStrTest = "postgres://postgres:1234@localhost:5432/test?sslmode=disable"

func TestDatabase(t *testing.T) {

	db, err := TestDB(connStrTest)
	if err != nil {
		t.Errorf("error connecting to database: %v", err)
	}
	query := `CREATE TABLE test_user (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL
	);
	CREATE TABLE test_squeak (
		id SERIAL PRIMARY KEY,
		user_id INT,
		squeak VARCHAR(255),
		FOREIGN KEY (user_id) REFERENCES test_user(id)
	)`

	_, err = db.Exec(query)
	if err != nil {
		t.Fatal(err)
	}

	store := PostgreSQLUserStore{db: db}

	t.Run("stores new squeak", func(t *testing.T) {
		user := "Mark"
		squeak := "I don't believe it!"

		err := store.PostSqueak(user, squeak)
		if err != nil {
			t.Error("error inserting data into database", err)
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
	_, err = db.Exec("DROP TABLE test_user, test_squeak")
	if err != nil {
		t.Fatalf("error dropping table: %v", err)
	}
}
