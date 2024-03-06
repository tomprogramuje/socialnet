package main

import (
	"log"
	"net/http"
)

func main() {
	db := NewPostgreSQLConnection(connStrProd)
	server := NewUserServer(NewPostgreSQLUserStore(db))
	log.Fatal(http.ListenAndServe(":8000", server))
}
