package main

import (
	"log"
	"net/http"

	"github.com/brangb/book_store_management/pkg/routes"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	routes.RegisterBookstoreRoutes(r)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":4000", r))

}
