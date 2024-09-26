package main

import (
	"log"
	"net/http"

	"github.com/brangb/golang_mongodb/controllers"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

func main() {
	// Initialize the router and controllers
	r := httprouter.New()
	uc := controllers.NewUserController(getSession())

	// Define routes
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)

	// Start the server
	log.Println("Server is running on port 4000...")
	log.Fatal(http.ListenAndServe(":4000", r))
}

// getSession establishes a MongoDB session
func getSession() *mgo.Session {
	// Connect to the MongoDB instance
	s, err := mgo.Dial("mongodb://localhost:27017")

	if err != nil {
		// log.Fatal("Failed to connect to MongoDB:", err)
		panic(err)
	}

	return s
}
