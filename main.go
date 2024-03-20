package main

import (
	"log"
	"net/http"
)

func main() {
	db, err := NewPostgreSQLConnection(connStrProd)
	if err != nil {
		log.Fatalf("problem securing connection to database: %v", err)
	}
	
	initializeDatabase(db)
	server := NewUserServer(NewPostgreSQLUserStore(db))
	log.Fatal(http.ListenAndServe(":8000", server))
}
