package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostingSqueaksAndRetrievingThem(t *testing.T) {
	store := NewInMemoryUserStore()
	server := UserServer{store}
	user := "Mark"

	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest(user))
	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest(user))
	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest(user))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetSqueakRequest(user))
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}
