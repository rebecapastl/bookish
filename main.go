package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var config *Config
var db *sql.DB
var err error

func main() {

	// configs
	config, err = LoadConfig()
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	// connect to db
	db, err = ConnectToDb(config.Database.URL)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	// create tables
	err = CreateTables(db)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the Books Database")
	})
	r.HandleFunc("/books", createBook).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book *Book
	var bookArgs BookArgs

	err := json.NewDecoder(r.Body).Decode(&bookArgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	book, err = CreateBook(db, bookArgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}