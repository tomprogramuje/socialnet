package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubUserStore struct {
	// Squeaks are Gopher's variant of tweets
	squeaks    map[string]int 
	newSqueaks []string 
}

func (s *StubUserStore) GetUserSqueakCount(name string) int {
	return s.squeaks[name]
}

func (s *StubUserStore) PostSqueak(name string) {
	s.newSqueaks = append(s.newSqueaks, name)
}

func TestStoreNewSqueaks(t *testing.T) {
	store := StubUserStore{
		map[string]int{},
		nil,
	}
	server := &UserServer{&store}

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
}

/*func TestPOSTSqueaks(t *testing.T) {
	store := StubUserStore{map[string]string{}}
	server := &UserServer{&store}

	t.Run("it saves a new squeak", func(t *testing.T) {
		request := newPostSqueakRequest("Andrew", "I am C-3PO, human-cyborg relations.")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if store.squeaks["Andrew"] != "I am C-3PO, human-cyborg relations." {
			t.Error("did not manage to save the right squeak")
		}
	})
}*/

func newPostSqueakRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s", name), nil)
	return req
}

func TestGETSqueaks(t *testing.T) {
	store := StubUserStore{
		map[string]int{
			"Mark":     12,
			"Harrison": 24,
		},
		nil,
	}
	server := &UserServer{&store}

	t.Run("returns Mark's squeak count", func(t *testing.T) {
		request := newGetSqueakRequest("Mark")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "12")
	})
	t.Run("returns Harrison's squeak count", func(t *testing.T) {
		request := newGetSqueakRequest("Harrison")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "24")
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
