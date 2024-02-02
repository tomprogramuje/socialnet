package main

import (
	"log"
	"net/http"
)

type InMemoryUserStore struct{}

func (i *InMemoryUserStore) GetUserPost(name string) string {
	return "hello everybody"
}

func main() {
	server := &UserServer{&InMemoryUserStore{}}
	log.Fatal(http.ListenAndServe(":80", server))
}
