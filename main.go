package main

import (
	"log"
	"net/http"
)

func main() {
	server := &UserServer{NewInMemoryUserStore()}
	log.Fatal(http.ListenAndServe(":80", server))
}
