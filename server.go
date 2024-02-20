package main

import (
	"encoding/json"
	"net/http"
)

type UserServer struct {
	store UserStore
	http.Handler
}

type User struct {
	Name    string
	Squeaks []string
}

func NewUserServer(store UserStore) *UserServer {
	u := new(UserServer)

	u.store = store

	router := http.NewServeMux()
	router.Handle("/userbase", http.HandlerFunc(u.userbaseHandler))
	router.Handle("GET /users/{name}", http.HandlerFunc(u.showSqueak))
	router.Handle("POST /users/{name}", http.HandlerFunc(u.saveSqueak))

	u.Handler = router

	return u
}

type UserStore interface {
	// Squeaks are Gopher's variant of tweets
	GetUserSqueaks(name string) []string
	PostSqueak(name, squeak string)
	GetUserbase() []User
}

const jsonContentType = "application/json"

func (u *UserServer) userbaseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(u.store.GetUserbase())
}

func (u *UserServer) showSqueak(w http.ResponseWriter, r *http.Request) {
	user := r.PathValue("name")
	squeaks := u.store.GetUserSqueaks(user)
	if len(squeaks) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", jsonContentType)
	if err := json.NewEncoder(w).Encode(squeaks); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (u *UserServer) saveSqueak(w http.ResponseWriter, r *http.Request) {
	user := r.PathValue("name")
	var payload User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Failed to decode JSON payload", http.StatusBadRequest)
		return
	}
	squeak := string(payload.Squeaks[0])
	u.store.PostSqueak(user, squeak)
	w.WriteHeader(http.StatusAccepted)
}
