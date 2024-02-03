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
	GetUserSqueak(name string) string
}

func (u *UserServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		u.saveSqueak(w)
	case http.MethodGet:
		u.showSqueak(w, r)
	}
}

func (u *UserServer) showSqueak(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimPrefix(r.URL.Path, "/users/")

	squeak := u.store.GetUserSqueak(user)

	if squeak == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, squeak)
}

func (u *UserServer) saveSqueak(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}
