package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostingSqueaksAndRetrievingThem(t *testing.T) {
	testStore := NewInMemoryUserStore()
	server := NewUserServer(testStore)
	user := User{"Mark", []string{"I don't believe it!"}}
	testStore.store["Mark"] = []string{"I don't believe it!"}

	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest(user.Name))
	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest(user.Name))
	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest(user.Name))

	t.Run("get squeak count", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetSqueakRequest(user.Name))

		got := getUserSqueaksFromResponse(t, response.Body)
		want := []string{"I don't believe it!"}

		assertResponse(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)

	})
	//t.Run("post to userbase", func(t *testing.T){})
	t.Run("get userbase", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newUserbaseRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getUserbaseFromResponse(t, response.Body)
		want := []User{
			{"Mark", []string{"I don't believe it!"}},
		}
		assertUserbase(t, got, want)
	})
}
