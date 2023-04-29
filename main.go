package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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
		r.HandleFunc("/books", CreateBookHandler).Methods("POST")
		r.HandleFunc("/books", ListBookHandler).Methods("GET")
		r.HandleFunc("/collections", ListCollectionHandler).Methods("GET")
		r.HandleFunc("/collections", CreateCollectionHandler).Methods("POST")
		r.HandleFunc("/collections/{collection_id}", AddBookToCollectionHandler).Methods("POST")


		log.Fatal(http.ListenAndServe(":8080", r))
	}
}

func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var book *Book
	var bookArgs BookArgs

	// decode request into arguments to function
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

	message := fmt.Sprintf("Book %s created with ID %d\n", book.Title, book.BookID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(message))
	json.NewEncoder(w).Encode(book)
}

func ListBookHandler(w http.ResponseWriter, r *http.Request) {
    bookArgs := &BookArgs{}

	// decode request into argumets to function
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

func CreateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	var collection *Collection
	var collectionArgs CollectionArgs

	// decode request into argumets to function
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
	
	message := fmt.Sprintf("Collection %s created with ID %d\n", collection.CollectionName, collection.CollectionID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(message))
	json.NewEncoder(w).Encode(collection)
}

func ListCollectionHandler(w http.ResponseWriter, r *http.Request) {
    collectionArgs := &CollectionArgs{}

	// decode request into argumets to function
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

func AddBookToCollectionHandler(w http.ResponseWriter, r *http.Request) {
	addArgs := &AddBookToCollectionArgs{}

	// decode request into arguments to function
	err = json.NewDecoder(r.Body).Decode(addArgs)
	if err != nil {
		if err.Error() == "EOF" {
			err = errors.New("no collection title set, collection not created")
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract collection_id from URL path
	vars := mux.Vars(r)
	collectionIDStr := vars["collection_id"]
	addArgs.CollectionID, err = SanitizeIdNumber(&collectionIDStr)
	if err != nil {
		http.Error(w, "Invalid collection ID", http.StatusBadRequest)
		return
	}

	// Call AddBookToCollection with the arguments
	collection, book, err := AddBookToCollection(db, *addArgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success status
	message := fmt.Sprintf("Book %s added to collection %s\n", book.Title, collection.CollectionName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
