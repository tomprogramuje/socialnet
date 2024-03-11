package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostingSqueaksAndRetrievingThem(t *testing.T) {
	db := NewPostgreSQLConnection(connStrTest)
	clearDatabase(db)
	initializeDatabase(db)
	testStore := NewPostgreSQLUserStore(db)
	server := NewUserServer(testStore)

	t.Run("create new user Harrison", func(t *testing.T) {
		body := []byte(`{
			"username": "Harrison", "email": "test", "password": "test"
		}`)
		
		request, _ := http.NewRequest(http.MethodPost, "/register/", bytes.NewBuffer(body))
		server.ServeHTTP(httptest.NewRecorder(), request)

		got, err := testStore.GetUserByName("Harrison")
		want := 1

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
		assertNoError(t, err)
	})
	t.Run("save squeaks for Harrison", func(t *testing.T) {
		response := httptest.NewRecorder()
		body := []byte(`
		{"name": "Harrison", "squeaks": ["Great, kid, don't get cocky."]}	
		`)
		server.ServeHTTP(response, newPostSqueakRequest("Harrison", body))

		body = []byte(`
			{"name": "Harrison", "squeaks": ["Laugh it up, fuzzball!"]}	
		`)
		server.ServeHTTP(response, newPostSqueakRequest("Harrison", body))

		assertStatus(t, response.Code, http.StatusAccepted)
	})
	t.Run("get Harrison's squeaks", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetSqueakRequest("Harrison"))

		got := getUserSqueaksFromResponse(t, response.Body)
		want := []string{"Great, kid, don't get cocky.", "Laugh it up, fuzzball!"}

		assertResponse(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("get userbase", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newUserbaseRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getUserbaseFromResponse(t, response.Body)
		want := []User{
			{"Harrison", "test", "test", []string{"Great, kid, don't get cocky.", "Laugh it up, fuzzball!"}},
		}

		assertUserbase(t, got, want)
	})
}
