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
	user := strings.TrimPrefix(r.URL.Path, "/users/")

	post := u.store.GetUserPost(user)

	if post == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, u.store.GetUserPost(user))
}
