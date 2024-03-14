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

type StubUserStore struct {
	userbase []User
}

func (s *StubUserStore) CreateUser(name, email, password string) (int, error) {
	id := len(s.userbase) + 1
	s.userbase = append(s.userbase, User{id, name, email, password, []string{}})
	return id, nil
}

func (s *StubUserStore) GetUserByUsername(username string) (*User, error) {
	var user *User
	for i := range s.userbase {
		if s.userbase[i].Username == username {
			user = &s.userbase[i]
			return user, nil
		}
	}
	return nil, fmt.Errorf("no user with that username (%s)found", username )
}

func (s *StubUserStore) PostSqueak(username, squeak string) (int, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return 0, fmt.Errorf("error getting user: %s", err)
	}
	user.Squeaks = append(user.Squeaks, squeak)
	return user.ID, nil
}

func (s *StubUserStore) GetUserSqueaks(name string) ([]string, error) {
	user, err := s.GetUserByUsername(name)
	if err != nil {
		return nil, err
	}
	return user.Squeaks, nil
}

func (s *StubUserStore) GetUserbase() ([]User, error) {
	return s.userbase, nil
}

func TestAuthentication(t *testing.T) {
	store := StubUserStore{}
	server := NewUserServer(&store)

	t.Run("returns correct userbase after registering new user", func(t *testing.T) {
		body := []byte(`{
			"username": "Carrie", 
			"email": "test",
			"password": "test"
		}`)
		request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		request = newUserbaseRequest()
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getUserbaseFromResponse(t, response.Body)
		want := []User{
			{1, "Carrie", "test", "", []string{}},
		}

		if len(got) != len(want) {
			t.Errorf("got %v users want %v users", len(got), len(want))
		}

		for i := range got {
			got[i].Password = ""
		}

		assertUserbase(t, got, want)
	})
	t.Run("password successfully verified", func(t *testing.T) {
		body := []byte(`{
			"username": "Carrie", 
			"password": "test"
		}`)
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
	t.Run("username and email already taken", func(t *testing.T) {

	})
	t.Run("user succesfully logged in", func(t *testing.T) {

	})
	t.Run("failed login", func(t *testing.T) {

	})
}

func TestStoreNewSqueaks(t *testing.T) {
	store := StubUserStore{
		[]User{{1, "Mark", "", "", []string{}}},
	}
	server := NewUserServer(&store)

	t.Run("it saves squeak on POST", func(t *testing.T) {
		body := []byte(`{
			"squeak": "Let go of your hate."
		}`)

		request := newPostSqueakRequest("Mark", body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.userbase[0].Squeaks) != 1 {
			t.Errorf("got %d calls to PostSqueak want %d", len(store.userbase[0].Squeaks), 1)
		}

		request = newGetSqueakRequest("Mark")
		server.ServeHTTP(response, request)

		got, err := store.GetUserSqueaks("Mark")

		assertResponse(t, got, []string{"Let go of your hate."})
		assertNoError(t, err)
	})
}

func newPostSqueakRequest(name string, body []byte) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s", name), bytes.NewBuffer(body))
	return req
}

func TestGETSqueaks(t *testing.T) {
	store := StubUserStore{
		[]User{
			{1, "Mark", "", "", []string{"I don't believe it!"}},
			{2, "Harrison", "", "", []string{"Great, kid, don't get cocky.", "Laugh it up, fuzzball!"}},
		},
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
			{1, "Mark", "", "", []string{"I don't believe it!"}},
			{2, "Harrison", "", "", []string{"I have a bad feeling about this.", "Great, kid, don't get cocky."}},
			{3, "Carrie", "", "", []string{"Will somebody get this big walking carpet out of my way?"}},
		}

		store := StubUserStore{wantedUserbase}
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
