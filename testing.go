package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"testing"
)

func newPostSqueakRequest(name string, body []byte) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s", name), bytes.NewBuffer(body))
	return req
}

func newGetSqueakRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", name), nil)
	return req
}

func newUserbaseRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/userbase", nil)
	return req
}

func getUserbaseFromResponse(t testing.TB, body io.Reader) (userbase []User) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&userbase)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of User, '%v'", body, err)
	}

	return
}

func getUserSqueaksFromResponse(t testing.TB, body io.Reader) (userSqueaks []string) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&userSqueaks)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of string, '%v'", body, err)
	}

	return
}

func assertResponse(t testing.TB, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("response is wrong, got %q, want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get the correct status, got %d, want %d", got, want)
	}
}

func assertUserbase(t testing.TB, got, want []User) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("userbases lenght mismatch: got %d want %d", len(got), len(want))
	}

	for i := range got {
		if got[i].ID != want[i].ID ||
			got[i].Username != want[i].Username ||
			got[i].Email != want[i].Email ||
			got[i].Password != want[i].Password ||
			!reflect.DeepEqual(got[i].Squeaks, want[i].Squeaks) {
			t.Errorf("user at index %d does not match: got %+v, want %+v", i, got[i], want[i])
		}
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func assertEqual[V comparable](t testing.TB, got, want V) {
	t.Helper()
	if got != want {
		t.Error("returned value differs from expected value, got", got, "want", want)
	}
}

func assertSqueaks(t testing.TB, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("did not get correct response, got %s, want %s", got, want)
	}
}

func assertError(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q back, but wanted %q", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}