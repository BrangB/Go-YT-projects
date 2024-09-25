package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allow", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Hello! You'r in the right path")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm Error: %v", err)
		return
	}

	fmt.Fprintf(w, "Post data successfully")
	name := r.FormValue("name")
	email := r.FormValue("email")

	fmt.Fprintf(w, "Your name is %s\n", name)
	fmt.Fprintf(w, "email is %s", email)

}

func main() {

	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileServer)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/form", formHandler)

	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal(err)
	}
}
