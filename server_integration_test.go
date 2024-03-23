package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPostingSqueaksAndRetrievingThem(t *testing.T) {
	db, err := NewPostgreSQLConnection(connStrTest)
	if err != nil {
		t.Fatal(err)
	}
	clearDatabase(db)
	initializeDatabase(db)
	testStore := NewPostgreSQLUserStore(db)
	server := NewUserServer(testStore)
	var jwtToken string

	t.Run("create new user Harrison", func(t *testing.T) {
		body := []byte(`{
			"username": "Harrison", "email": "test", "password": "test"
		}`)

		request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		server.ServeHTTP(httptest.NewRecorder(), request)

		got, err := testStore.GetUserByUsername("Harrison")
		want := 1

		if got.ID != want {
			t.Errorf("got %d want %d", got.ID, want)
		}
		assertNoError(t, err)
	})
	t.Run("save squeaks for Harrison without logging in", func(t *testing.T) {
		response := httptest.NewRecorder()
		body := []byte(`
		{"text": "Great, kid, don't get cocky."}	
		`)
		server.ServeHTTP(response, newPostSqueakRequest("Harrison", body))

		body = []byte(`
			{"text": "Laugh it up, fuzzball!"}	
		`)
		server.ServeHTTP(response, newPostSqueakRequest("Harrison", body))

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})
	t.Run("logging in as Harrison", func(t *testing.T) {
		body := []byte(`{
			"username": "Harrison",
			"password": "test"
		}`)
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		for _, cookie := range response.Result().Cookies() {
			if cookie.Name == "Authorization" {
				jwtToken = cookie.Value
				break
			}
		}

		if jwtToken == "" {
			t.Error("JWT token not found in response cookies")
		} 
	})
	t.Run("save squeaks for Harrison after successful login", func(t *testing.T) {
		body := []byte(`
			{"text": "Great, kid, don't get cocky."}	
		`)
		request := newPostSqueakRequestWithJWT("Harrison", body, jwtToken)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		body = []byte(`
			{"text": "Laugh it up, fuzzball!"}	
		`)
		request = newPostSqueakRequestWithJWT("Harrison", body, jwtToken)
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
	t.Run("get Harrison's squeaks", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetSqueakRequest("Harrison"))

		got := getUserSqueaksFromResponse(t, response.Body)
		want := []SqueakPost{{"Great, kid, don't get cocky.", time.Now()}, {"Laugh it up, fuzzball!", time.Now()}}

		assertResponse(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("get userbase", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newUserbaseRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getUserbaseFromResponse(t, response.Body)
		want := []User{
			{1, "Harrison", "test", "", []SqueakPost{{"Great, kid, don't get cocky.", time.Now()}, {"Laugh it up, fuzzball!", time.Now()}}, time.Now()},
		}

		if len(got) != len(want) {
			t.Errorf("got %v users want %v users", len(got), len(want))
		}

		for i := range got {
			got[i].Password = ""
		}

		assertUserbase(t, got, want)
	})
}
