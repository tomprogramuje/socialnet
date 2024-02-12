package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"testing"
)

type StubUserStore struct {
	// Squeaks are Gopher's variant of tweets
	squeaks    map[string][]string
	newSqueaks []string
	userbase   []User
}

func (s *StubUserStore) GetUserSqueaks(name string) []string {
	return s.squeaks[name]
}

func (s *StubUserStore) PostSqueak(name string) {
	s.newSqueaks = append(s.newSqueaks, name)
}

func (s *StubUserStore) GetUserbase() []User {
	return s.userbase
}

func TestStoreNewSqueaks(t *testing.T) {
	store := StubUserStore{
		map[string][]string{},
		nil,
		nil,
	}
	server := NewUserServer(&store)

	t.Run("it records squeaks on POST", func(t *testing.T) {
		user := "Mark"

		request := newPostSqueakRequest(user)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.newSqueaks) != 1 {
			t.Errorf("got %d calls to RecordSqueak want %d", len(store.newSqueaks), 1)
		}

		if store.newSqueaks[0] != user {
			t.Errorf("did not store correct user got %q want %q", store.newSqueaks[0], user)
		}
	})
	//t.Run("post to userbase", func(t *testing.T){}) - the previous test should probably be just part of this or better yet all storing tests should be part of db_test kit 
}

func newPostSqueakRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s", name), nil)
	return req
}

func TestGETSqueaks(t *testing.T) {
	store := StubUserStore{
		map[string][]string{
			"Mark":     {"I don't believe it!"},
			"Harrison": {"Great, kid, don't get cocky.", "Laugh it up, fuzzball!"},
		},
		nil,
		nil,
	}
	server := NewUserServer(&store)

	t.Run("returns Mark's squeak", func(t *testing.T) {
		request := newGetSqueakRequest("Mark")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getUserSqueaksFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponse(t, got, []string{"I don't believe it!"})
		assertContentType(t, response, jsonContentType)
	})
	t.Run("returns Harrison's squeaks", func(t *testing.T) {
		request := newGetSqueakRequest("Harrison")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getUserSqueaksFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponse(t, got, []string{"Great, kid, don't get cocky.", "Laugh it up, fuzzball!"})
		assertContentType(t, response, jsonContentType)
	})
	t.Run("returns 404 on missing user", func(t *testing.T) {
		request := newGetSqueakRequest("Carrie")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestUserbase(t *testing.T) {

	t.Run("it returns the user base as JSON", func(t *testing.T) {
		wantedUserbase := []User{
			{"Mark", []string{"I don't believe it!"}},
			{"Harrison", []string{"I have a bad feeling about this.", "Great, kid, don't get cocky."}},
			{"Carrie", []string{"Will somebody get this big walking carpet out of my way?"}},
		}

		store := StubUserStore{nil, nil, wantedUserbase}
		server := NewUserServer(&store)

		request := newUserbaseRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getUserbaseFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertUserbase(t, got, wantedUserbase)
		assertContentType(t, response, jsonContentType)
	})
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}
