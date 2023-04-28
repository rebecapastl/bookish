package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

func CreateAuthor(db *sql.DB, d AuthorArgs) (*Author, error){
	var author Author
	var err error

	err = db.QueryRow("INSERT INTO authors (name) VALUES ($1) ON CONFLICT DO NOTHING RETURNING author_id, name, creation_date", d.Name).Scan(&author.AuthorID, &author.Name, &author.CreationDate)
    if err != nil {
        if err == sql.ErrNoRows{
            err = errors.New("author already exists in the database")
        }
        return nil, err
    }
	fmt.Printf("Author %s created with ID %d\n", author.Name, author.AuthorID)

	return &author, nil
}

func ListAuthors(db *sql.DB, d AuthorArgs) ([]Author, error) {
	var authors []Author

	query := "SELECT * FROM authors"
	if d.Name != nil {
		query += " WHERE name = '" + *d.Name + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var author Author
		err := rows.Scan(&author.AuthorID, &author.Name, &author.CreationDate)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func CreateBook(db *sql.DB, m BookArgs) (*Book, error){
    var author *Author
    var err error

    // try to fetch said author
    authors, err := ListAuthors(db, AuthorArgs{Name: m.Author})
	if err != nil{
		return nil, err
	}

    // create a new one in case the author is not in the db
    if len(authors) == 0{
        author, err = CreateAuthor(db, AuthorArgs{Name: m.Author})
        if err != nil{
            return nil, err
        }
    } else {
        author = &authors[0]
    }

    var book Book
    err = db.QueryRow("INSERT INTO books (title, author_id) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING book_id, title, creation_date", m.Title, author.AuthorID).Scan(&book.BookID, &book.Title, &book.CreationDate)
    if err != nil {
        if err == sql.ErrNoRows{
            err = errors.New("book already exists in the database")
        }
        return nil, err
    }

    book.Author = author.Name
    fmt.Printf("Book %s created with ID %d\n", book.Title, book.BookID)

    return &book, nil
}

func ListBooks(db *sql.DB, m BookArgs) ([]Book, error) {
    var books []Book

    query := `
        SELECT books.book_id, books.title, authors.name, books.creation_date
        FROM books
        JOIN authors ON books.author_id = authors.author_id
        `

    whereClauses := []string{}
    if m.BookID != nil {
        whereClauses = append(whereClauses, fmt.Sprintf("books.book_id = %d", *m.BookID))
    }
    if m.Title != nil {
        whereClauses = append(whereClauses, fmt.Sprintf("books.title = '%s'", *m.Title))
    }
    if m.Author != nil {
        whereClauses = append(whereClauses, fmt.Sprintf("authors.name = '%s'", *m.Author))
    }

    if len(whereClauses) > 0 {
        query += "WHERE " + strings.Join(whereClauses, " AND ")
    }

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var book Book
        err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.CreationDate)
        if err != nil {
            return nil, err
        }
        books = append(books, book)
    }

	err = rows.Err()
    if err != nil {
        return nil, err
    }

	if len(books) == 0{
		err = errors.New("no books with the chosen specification")
        return nil, err
	}

    return books, nil
}

func CreateCollection(db *sql.DB, c CollectionArgs) (*Collection, error){
	var collection Collection

	err := db.QueryRow("INSERT INTO collections (collection_name) VALUES ($1) ON CONFLICT DO NOTHING RETURNING collection_id, collection_name, creation_date", c.CollectionName).Scan(&collection.CollectionID, &collection.CollectionName, &collection.CreationDate)
	if err != nil {
        if err == sql.ErrNoRows{
            err = errors.New("collection already exists in the database")
        }
        return nil, err
	}
	fmt.Printf("Collection %s created with ID %d\n", collection.CollectionName, collection.CollectionID)

	return &collection, nil
}

func ListCollections(db *sql.DB, c CollectionArgs) ([]Collection, error) {
    var collections []Collection

    query := "SELECT * FROM collections"

    whereClauses := []string{}
    if c.CollectionID != nil {
        whereClauses = append(whereClauses, fmt.Sprintf("collection_id = %d", *c.CollectionID))
    }
    if c.CollectionName != nil {
        whereClauses = append(whereClauses, fmt.Sprintf("collection_name = '%s'", *c.CollectionName))
    }

    if len(whereClauses) > 0 {
        query += "WHERE " + strings.Join(whereClauses, " AND ")
    }

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    err = rows.Err()
    if err != nil {
        return nil, err
    }

	if len(collections) == 0{
        err = errors.New("no collections with the chosen specification")
		return nil, err
	}

    return collections, nil
}
