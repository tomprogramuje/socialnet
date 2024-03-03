package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewUserServer(NewInMemoryUserStore())
	log.Fatal(http.ListenAndServe(":8000", server))
}

// todo: use PostgresqlUserStore