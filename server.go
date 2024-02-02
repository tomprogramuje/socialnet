package main

import (
	"fmt"
	"net/http"
	"strings"
)

func PostServer(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimPrefix(r.URL.Path, "/users/")

	fmt.Fprint(w, GetUserPost(user))
}

func GetUserPost(name string) string {
	if name == "Mark" {
		return "Hey, how is everybody today?"
	}

	if name == "Harrison" {
		return "I am having an awful day..."
	}

	return ""
}
