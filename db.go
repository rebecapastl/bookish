package main

import (
	"database/sql"
	"fmt"
	
	_ "github.com/lib/pq"
)

func ConnectToDb(url string) (*sql.DB, error) {

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	} else {
		fmt.Println("openned database")
	}

	if err := db.Ping(); err != nil {
		return nil, err
	} else {
		fmt.Println("connected to database")
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	// TODO: modify database to accept multiple authors in one book (add table BOOK_AUTHOR to represent this relationship, remove author_id from BOOKS, change methods accordingly)

	// create authors table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS authors (
        author_id SERIAL PRIMARY KEY,
        name VARCHAR(100) UNIQUE NOT NULL, CHECK (name <> ''),
		creation_date DATE DEFAULT CURRENT_DATE
    );`)
	if err != nil {
		return err
	} else {
		fmt.Println("table authors ok")
	}

	// create collections table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS collections (
        collection_id SERIAL PRIMARY KEY,
        collection_name VARCHAR(50) UNIQUE NOT NULL, CHECK (collection_name <> ''),
        creation_date DATE DEFAULT CURRENT_DATE
    );`)
	if err != nil {
		return err
	} else {
		fmt.Println("table collections ok")
	}

	// create books table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS books (
        book_id SERIAL PRIMARY KEY,
        title VARCHAR(100) NOT NULL, CHECK (title <> ''),
		release_date DATE,
		edition_number INT,
		creation_date DATE DEFAULT CURRENT_DATE,
        author_id INT NOT NULL,
        FOREIGN KEY (author_id) REFERENCES authors(author_id),
		UNIQUE (title, author_id)
    );`)
	if err != nil {
		return err
	} else {
		fmt.Println("table books ok")
	}

	// create book_in_collection table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS book_in_collection (
		book_id INT,
		collection_id INT,
		FOREIGN KEY (book_id) REFERENCES books(book_id) ON DELETE CASCADE,
		FOREIGN KEY (collection_id) REFERENCES collections(collection_id) ON DELETE CASCADE,
		PRIMARY KEY (book_id, collection_id)
    );`)
	if err != nil {
		return err
	} else {
		fmt.Println("table book_in_collection ok")
	}

	return nil
}