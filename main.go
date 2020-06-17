package main

import (
	"fmt"
	"log"
	"net/http"
)

var inMemoryDb *InMemoryStore

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		command := r.FormValue("command")
		result := inMemoryDb.ProcessCommand(command)
		fmt.Fprintf(w, "%s\n", result)

	default:
		fmt.Fprintf(w, "Sorry, only POST method supported.")
	}
}

func main() {
	inMemoryDb = CreateInMemStore()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
