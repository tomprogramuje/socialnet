package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostingSqueaksAndRetrievingThem(t *testing.T) {
	testStore := NewInMemoryUserStore()
	server := NewUserServer(testStore)
	bodyMark := []byte(`
		{"name": "Mark", "squeaks": ["I don't believe it!"]}
	`)

	bodyHarrison := []byte(`
		{"name": "Harrison", "squeaks": ["Great, kid, don't get cocky."]}	
	`)

	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest("Mark", bodyMark))
	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest("Harrison", bodyHarrison))

	bodyHarrison = []byte(`
		{"name": "Harrison", "squeaks": ["Laugh it up, fuzzball!"]}	
	`)

	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest("Harrison", bodyHarrison))

	t.Run("get Mark's squeaks", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetSqueakRequest("Mark"))

		got := getUserSqueaksFromResponse(t, response.Body)
		want := []string{"I don't believe it!"}

		assertResponse(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)

	})
	
	t.Run("get userbase", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newUserbaseRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getUserbaseFromResponse(t, response.Body)
		want := []User{
			{"Mark", []string{"I don't believe it!"}},
			{"Harrison", []string{"Great, kid, don't get cocky.", "Laugh it up, fuzzball!"}},
		}
		assertUserbase(t, got, want)
	})
}
