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
	inMemoryDb = CreateInMemStore(5, "AOF.log")
	http.HandleFunc("/", handler)
	fmt.Println("Server starting at at http://localhost:8080/ use ctrl+c to stop.\n" +
		"You can send commads as x-www-form-urlencoded POST request key value eg. 'command=SET k1 v1' \n" +
		"Eg.:\n\ncurl -d 'command=SET edtech=awesome' http://localhost:8080/\n\n ")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
