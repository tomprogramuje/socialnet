package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPosts(t *testing.T) {
	t.Run("returns Mark's post", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/users/Mark", nil)
		response := httptest.NewRecorder()

		PostServer(response, request)

		got := response.Body.String()
		want := "Hey, how is everybody today?"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("returns Harrison's post", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/users/Harrison", nil)
		response := httptest.NewRecorder()

		PostServer(response, request)

		got := response.Body.String()
		want := "I am having an awful day..."

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
