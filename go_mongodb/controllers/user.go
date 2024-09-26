package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brangb/golang_mongodb/models"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

// GetUser retrieves a user by their ID
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	session := uc.session.Clone()
	defer session.Close()
	id := p.ByName("id")

	// Check if the provided ID is a valid BSON ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(id)
	u := models.User{}

	// Retrieve the user from the database
	if err := uc.session.DB("mongo-golang").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error marshalling user:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

// CreateUser creates a new user
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := uc.session.Clone()
	defer session.Close()
	u := models.User{} // Create an instance of the User model

	// Decode the request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error decoding user:", err)
		return
	}

	// Set a new ObjectId for the user
	u.Id = bson.NewObjectId()

	// Insert the user into the database
	if err := uc.session.DB("mongo-golang").C("users").Insert(u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error inserting user:", err)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error marshalling user:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

// DeleteUser deletes a user by their ID
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	session := uc.session.Clone()
	defer session.Close()
	id := p.ByName("id")

	// Check if the provided ID is a valid BSON ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(id)

	// Remove the user from the database
	if err := uc.session.DB("mongo-golang").C("users").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user %s\n", oid.Hex())
}
