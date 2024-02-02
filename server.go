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
	GetUserPost(name string) string
}

func (u *UserServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		u.savePost(w)
	case http.MethodGet:
		u.showPost(w, r)
	}
}

func (u *UserServer) showPost(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimPrefix(r.URL.Path, "/users/")

	post := u.store.GetUserPost(user)

	if post == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, post)
}

func (u *UserServer) savePost(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}
