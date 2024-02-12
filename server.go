package main

import (
	"encoding/json"
	"net/http"
	"strings"
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
	router.Handle("/users/", http.HandlerFunc(u.usersHandler))

	u.Handler = router

	return u
}

type UserStore interface {
	// Squeaks are Gopher's variant of tweets
	GetUserSqueaks(name string) []string
	PostSqueak(name string)
	GetUserbase() []User
}

const jsonContentType = "application/json"

func (u *UserServer) userbaseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(u.store.GetUserbase())
}

func (u *UserServer) usersHandler(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimPrefix(r.URL.Path, "/users/")

	switch r.Method {
	case http.MethodPost:
		u.saveSqueak(w, user) // I need to send the squeak somehow, JSON?
	case http.MethodGet:
		u.showSqueak(w, user)
	}
}

func (u *UserServer) showSqueak(w http.ResponseWriter, user string) {
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

func (u *UserServer) saveSqueak(w http.ResponseWriter, user string) {
	u.store.PostSqueak(user)
	w.WriteHeader(http.StatusAccepted)
}
