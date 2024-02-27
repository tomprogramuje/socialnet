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
	query := `CREATE TABLE IF NOT EXISTS "user" (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS "squeak" (
		id SERIAL PRIMARY KEY,
		user_id INT,
		text VARCHAR(255),
		FOREIGN KEY (user_id) REFERENCES "user"(id)
	)`

	_, err = db.Exec(query)
	if err != nil {
		t.Fatal(err)
	}

	store := PostgreSQLUserStore{db: db}

	t.Run("creates new User", func(t *testing.T) {
		name := "Mark"

		got := store.CreateUser(db, name)
		if err != nil {
			t.Errorf("error creating new User %v", err)
		}

		want := 1
		if got != want {
			t.Errorf("got wrong id back, got %d want %d", got, want)
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
	_, err = db.Exec(`DROP TABLE squeak; DROP TABLE "user";`)
	if err != nil {
		t.Fatalf("error dropping table: %v", err)
	}
}
