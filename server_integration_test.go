package main

/*import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostingSqueaksAndRetrievingThem(t *testing.T) {
	store := NewInMemoryUserStore()
	server := NewUserServer(store)
	user := User{"Mark", []string{"I don't believe it!"}}

	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest(user.Name))
	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest(user.Name))
	server.ServeHTTP(httptest.NewRecorder(), newPostSqueakRequest(user.Name))

	t.Run("get squeak count", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetSqueakRequest(user.Name))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponse(t, response.Body.String(), "3")
	})
	//t.Run("post to userbase", func(t *testing.T){})
	t.Run("get userbase", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newUserbaseRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getUserbaseFromResponse(t, response.Body)
		want := []User{
			{"Mark", []string{"I don't believe it!"}},
		}
		assertUserbase(t, got, want)
	})
}
*/