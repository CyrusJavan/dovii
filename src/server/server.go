package server

import (
	"fmt"
	"log"
	"net/http"
)

// StartServer starts the server
func StartServer() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":7070", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
