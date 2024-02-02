package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPosts(t *testing.T) {
	t.Run("returns Mark's post", func(t *testing.T) {
		request := newGetPostRequest("Mark")
		response := httptest.NewRecorder()

		PostServer(response, request)

		assertResponseBody(t, response.Body.String(), "Hey, how is everybody today?")
	})
	t.Run("returns Harrison's post", func(t *testing.T) {
		request := newGetPostRequest("Harrison")
		response := httptest.NewRecorder()

		PostServer(response, request)

		assertResponseBody(t, response.Body.String(), "I am having an awful day...")
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
