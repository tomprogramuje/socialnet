package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := NewPostgreSQLConnection(os.Getenv("CONN_STR_PROD"))
	if err != nil {
		log.Fatalf("problem securing connection to database: %v", err)
	}

	initializeDatabase(db)
	server := NewUserServer(NewPostgreSQLUserStore(db))
	log.Fatal(http.ListenAndServe(":8000", server))
}
