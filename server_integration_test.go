package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostingSqueaksAndRetrievingThem(t *testing.T) {
	testStore := NewInMemoryUserStore()
	server := NewUserServer(testStore)
	body := []byte(`
		{"name": "Harrison", "squeaks": ["Great, kid, don't get cocky."]}	
	`)

	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest("Harrison", body))

	body = []byte(`
		{"name": "Harrison", "squeaks": ["Laugh it up, fuzzball!"]}	
	`)

	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest("Harrison", body))

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
			{"Harrison", []string{"Great, kid, don't get cocky.", "Laugh it up, fuzzball!"}},
		}
		assertUserbase(t, got, want)
	})
}
