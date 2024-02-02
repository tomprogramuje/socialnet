package main

import (
	"fmt"
	"net/http"
	"strings"
)

func PostServer(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimPrefix(r.URL.Path, "/users/")

	if user == "Mark" {
		fmt.Fprint(w, "Hey, how is everybody today?")
		return
	}

	if user == "Harrison" {
		fmt.Fprint(w, "I am having an awful day...")
		return
	}
}
