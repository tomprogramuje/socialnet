package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type StubUserStore struct {
	userbase []User
}

func (s *StubUserStore) CreateUser(name, email, password string) (int, error) {
	id := len(s.userbase) + 1
	s.userbase = append(s.userbase, User{id, name, email, password, []SqueakPost{}, time.Now()})
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
	return nil, fmt.Errorf("no user with that username (%s) found", username)
}

func (s *StubUserStore) PostSqueak(username, squeak string) (int, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return 0, fmt.Errorf("error getting user: %s", err)
	}
	user.Squeaks = append(user.Squeaks, SqueakPost{squeak, time.Now()})
	return user.ID, nil
}

func (s *StubUserStore) GetUserSqueaks(name string) ([]SqueakPost, error) {
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
			{1, "Carrie", "test", "", []SqueakPost{}, time.Now()},
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
		body := []byte(`{
			"username": "Carrie",
			"email": "test",
			"password": "test"
		}`)

		request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusConflict)
		assertError(t, response.Body.String(), `username already taken
email already taken
`)
	})
	t.Run("user succesfully logged in", func(t *testing.T) {

	})
	t.Run("failed login", func(t *testing.T) {

	})
}

func TestStoreNewSqueaks(t *testing.T) {
	store := StubUserStore{
		[]User{{1, "Mark", "", "", []SqueakPost{}, time.Now()}},
	}
	server := NewUserServer(&store)

	t.Run("it saves squeak on POST", func(t *testing.T) {
		body := []byte(`{
			"text": "Let go of your hate."
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

		assertResponse(t, got, []SqueakPost{{"Let go of your hate.", time.Now()}})
		assertNoError(t, err)
	})
}

func TestGETSqueaks(t *testing.T) {
	store := StubUserStore{
		[]User{
			{1, "Mark", "", "", []SqueakPost{{"I don't believe it!", time.Now()}}, time.Now()},
			{2, "Harrison", "", "", []SqueakPost{{"Great, kid, don't get cocky.", time.Now()}, {"Laugh it up, fuzzball!", time.Now()}}, time.Now()},
		},
	}
	server := NewUserServer(&store)

	t.Run("returns Mark's squeak", func(t *testing.T) {
		request := newGetSqueakRequest("Mark")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getUserSqueaksFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponse(t, got, []SqueakPost{{"I don't believe it!", time.Now()}})
		assertContentType(t, response, jsonContentType)
	})
	t.Run("returns Harrison's squeaks", func(t *testing.T) {
		request := newGetSqueakRequest("Harrison")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getUserSqueaksFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponse(t, got, []SqueakPost{{"Great, kid, don't get cocky.", time.Now()}, {"Laugh it up, fuzzball!", time.Now()}})
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
			{1, "Mark", "", "", []SqueakPost{{"I don't believe it!", time.Now()}}, time.Now()},
			{2, "Harrison", "", "", []SqueakPost{{"I have a bad feeling about this.", time.Now()}, {"Great, kid, don't get cocky.", time.Now()}}, time.Now()},
			{3, "Carrie", "", "", []SqueakPost{{"Will somebody get this big walking carpet out of my way?", time.Now()}}, time.Now()},
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
