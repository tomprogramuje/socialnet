package main

import (
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(PostServer)
	log.Fatal(":80", handler)
}
