package server

import (
	"fmt"
	"net/http"
)

func PostServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hey, how is everybody today?")
}
