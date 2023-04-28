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

	if len(os.Args) > 1 {
		CLIcommands()
	} else {
		r := mux.NewRouter()
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Welcome to the Books Database")
		})
		r.HandleFunc("/books", createBook).Methods("POST")
		r.HandleFunc("/books", listBooks).Methods("GET")
		r.HandleFunc("/collections", listCollections).Methods("GET")
		r.HandleFunc("/collections", createCollection).Methods("POST")
		r.HandleFunc("/collections/add_book", addBookToCollection).Methods("POST")


		log.Fatal(http.ListenAndServe(":8080", r))
	}
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

func listBooks(w http.ResponseWriter, r *http.Request) {
    bookArgs := &BookArgs{}
	if r.ContentLength != 0 {
		err = json.NewDecoder(r.Body).Decode(bookArgs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			
			return
		}
	}

    books, err := ListBooks(db, *bookArgs)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(books)
}

func createCollection(w http.ResponseWriter, r *http.Request) {
	var collection *Collection
	var collectionArgs CollectionArgs

	err := json.NewDecoder(r.Body).Decode(&collectionArgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	collection, err = CreateCollection(db, collectionArgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(collection)
}

func listCollections(w http.ResponseWriter, r *http.Request) {
    collectionArgs := &CollectionArgs{}
	if r.ContentLength != 0 {
		err = json.NewDecoder(r.Body).Decode(collectionArgs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			
			return
		}
	}

    collections, err := ListCollections(db, *collectionArgs)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(collections)
}

func addBookToCollection(w http.ResponseWriter, r *http.Request) {
	addArgs := &AddBookToCollectionArgs{}
    err := json.NewDecoder(r.Body).Decode(&addArgs)
    if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	collection, book, err := AddBookToCollection(db, *addArgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Book %s added to collection %s\n", book.Title, collection.CollectionName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
    
}
