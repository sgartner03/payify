package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var items ItemCouchbase

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the items API!")
	fmt.Println("Endpoint: Homepage")
}

func allItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: All Items")
	vars := mux.Vars(r)
	json.NewEncoder(w).Encode(items.GetAll(vars["username"]))
}

func singleItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	user := vars["username"]
	id, err := strconv.Atoi(key)
	if err != nil {
		log.Fatal(err)
	}

	item := items.Get(id, user)
	json.NewEncoder(w).Encode(item)
}

func createNewItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["username"]
	reqBody, _ := ioutil.ReadAll(r.Body)

	var item Item
	json.Unmarshal(reqBody, &item)
	item.Username = user
	items.Add(item)
	json.NewEncoder(w).Encode(item)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/{username}/all", allItems)
	router.HandleFunc("/{username}/item", createNewItem).Methods("POST")
	router.HandleFunc("/{username}/item/{id}", singleItem)

	log.Fatal(http.ListenAndServe(":8082", router))
}

func main() {
	items = NewCouchbase()
	handleRequests()
}
