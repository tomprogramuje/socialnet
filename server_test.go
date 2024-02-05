package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubUserStore struct {
	// Squeaks are Gopher's variant of tweets
	squeaks map[string]string
}

func (s *StubUserStore) GetUserSqueak(name string) string {
	squeak := s.squeaks[name]
	return squeak
}

func TestPOSTSqueaks(t *testing.T) {
	store := StubUserStore{map[string]string{}}
	server := &UserServer{&store}

	t.Run("it returns accepted on POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users/Mark", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
}

func TestGETSqueaks(t *testing.T) {
	store := StubUserStore{
		map[string]string{
			"Mark":     "Hey, how is everybody today?",
			"Harrison": "I am having an awful day...",
		},
	}
	server := &UserServer{&store}

	t.Run("returns Mark's squeak", func(t *testing.T) {
		request := newGetSqueakRequest("Mark")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "Hey, how is everybody today?")
	})
	t.Run("returns Harrison's squeak", func(t *testing.T) {
		request := newGetSqueakRequest("Harrison")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "I am having an awful day...")
	})
	t.Run("returns 404 on missing user", func(t *testing.T) {
		request := newGetSqueakRequest("Carrie")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
	t.Run("returns Andrew's squeaks", func(t *testing.T) {

	})
}

func newGetSqueakRequest(name string) *http.Request {
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
