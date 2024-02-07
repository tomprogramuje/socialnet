package main

import (
	"fmt"
	"net/http"
	"strings"
)

type UserServer struct {
	store UserStore
}

type UserStore interface {
	// Squeaks are Gopher's variant of tweets
	GetUserSqueakCount(name string) int
	PostSqueak(name string)
}

func (u *UserServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimPrefix(r.URL.Path, "/users/")

	switch r.Method {
	case http.MethodPost:
		u.saveSqueak(w, user) // I need to send the squeak somehow, JSON?
	case http.MethodGet:
		u.showSqueak(w, user)
	}
}

func (u *UserServer) showSqueak(w http.ResponseWriter, user string) {
	squeak := u.store.GetUserSqueakCount(user)

	if squeak == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, squeak)
}

func (u *UserServer) saveSqueak(w http.ResponseWriter, user string) {
	u.store.PostSqueak(user)
	w.WriteHeader(http.StatusAccepted)
}
