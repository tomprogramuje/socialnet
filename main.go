package main

import (
	"log"
	"net/http"
)

func main() {
	db := NewPostgreSQLConnection(connStrProd)
	initializeDatabase(db)
	server := NewUserServer(NewPostgreSQLUserStore(db))
	log.Fatal(http.ListenAndServe(":8000", server))
}
