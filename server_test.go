package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubUserStore struct {
	posts map[string]string
}

func (s *StubUserStore) GetUserPost(name string) string {
	post := s.posts[name]
	return post
}

func TestGETPosts(t *testing.T) {
	store := StubUserStore{
		map[string]string{
			"Mark":     "Hey, how is everybody today?",
			"Harrison": "I am having an awful day...",
		},
	}
	server := &UserServer{&store}

	t.Run("returns Mark's post", func(t *testing.T) {
		request := newGetPostRequest("Mark")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "Hey, how is everybody today?")
	})
	t.Run("returns Harrison's post", func(t *testing.T) {
		request := newGetPostRequest("Harrison")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "I am having an awful day...")
	})
	t.Run("returns 404 on missing user", func(t *testing.T) {
		request := newGetPostRequest("Carrie")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)	
	})
}

func newGetPostRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q, want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get the correct status, got %d, want %d", got, want)
	}
}