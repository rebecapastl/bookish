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
		fmt.Println("opened database")
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
	}

	// create collections table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS collections (
        collection_id SERIAL PRIMARY KEY,
        collection_name VARCHAR(50) UNIQUE NOT NULL, CHECK (collection_name <> ''),
        creation_date DATE DEFAULT CURRENT_DATE
    );`)
	if err != nil {
		return err
	}

	// create books table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS books (
        book_id SERIAL PRIMARY KEY,
        title VARCHAR(100) NOT NULL, CHECK (title <> ''),
		published_date DATE,
		edition_number INT,
		creation_date DATE DEFAULT CURRENT_DATE,
        author_id INT NOT NULL,
        FOREIGN KEY (author_id) REFERENCES authors(author_id),
		UNIQUE (title, author_id)
    );`)
	if err != nil {
		return err
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
	}

	return nil
}

func CreateAuthor(db *sql.DB, a AuthorArgs) (*Author, error){
	var author Author
	var err error

    // make sure author is not null(empty)
	a.Name = SanitizeAuthorName(a.Name)

	err = db.QueryRow("INSERT INTO authors (name) VALUES ($1) ON CONFLICT DO NOTHING RETURNING author_id, name, creation_date", a.Name).Scan(&author.AuthorID, &author.Name, &author.CreationDate)
    if err != nil {
        if err == sql.ErrNoRows{
            err = errors.New("author already exists in the database")
        }
        return nil, err
    }
	fmt.Printf("Author %s created with ID %d\n", author.Name, author.AuthorID)

	return &author, nil
}

func ListAuthors(db *sql.DB, a AuthorArgs) ([]Author, error) {
	var authors []Author

	query := "SELECT * FROM authors"
	if a.Name != nil {
		query += " WHERE name = '" + *a.Name + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() { // Use Next to advance from row to row. It prepares the next result row for reading with the Scan, even the first call must be preceded by a call to Next.
		var author Author
		err := rows.Scan(&author.AuthorID, &author.Name, &author.CreationDate)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

    err = rows.Err()
	if err != nil {
		return nil, err
	}

	return authors, nil
}

func CreateBook(db *sql.DB, b BookArgs) (*Book, error){
    var author *Author
    var err error

    if b.Title == nil || *b.Title == "" {
        return nil, errors.New("no book title set, book not created") 
    }

    // try to fetch said author, if there is no author no error returns
    // it is just a heads up to the CreateBook function, so it creates a new author for the book
    authors, err := ListAuthors(db, AuthorArgs{Name: b.Author})
	if err != nil{
		return nil, err
	}

    // create a new one in case the author is not in the db
    if len(authors) == 0{
        author, err = CreateAuthor(db, AuthorArgs{Name: b.Author})
        if err != nil{
            return nil, err
        }
    } else {
        author = &authors[0]
    }

    var book Book
    err = db.QueryRow("INSERT INTO books (title, author_id) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING book_id, title, creation_date", b.Title, author.AuthorID).Scan(&book.BookID, &book.Title, &book.CreationDate)
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

func ListBooks(db *sql.DB, b BookArgs) ([]Book, error) {
    var books []Book

    query := `
        SELECT books.book_id, books.title, authors.name, books.creation_date
        FROM books
        JOIN authors ON books.author_id = authors.author_id
        `

	// add to the wehre clause if it was present in the request args
    whereClauses := []string{}
    if b.BookID != nil {
        whereClauses = append(whereClauses, fmt.Sprintf("books.book_id = %d", *b.BookID))
    }
    if b.Title != nil {
        whereClauses = append(whereClauses, fmt.Sprintf("books.title = '%s'", *b.Title))
    }
    if b.Author != nil {
        whereClauses = append(whereClauses, fmt.Sprintf("authors.name = '%s'", *b.Author))
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

    if c.CollectionName == nil || *c.CollectionName == "" {
        return nil, errors.New("no collection name set, collection not created") 
    }

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

    query := `
        SELECT collections.collection_id, collections.collection_name, collections.creation_date,
               books.book_id, books.title, books.creation_date, authors.name
        FROM collections
        LEFT JOIN book_in_collection ON collections.collection_id = book_in_collection.collection_id
        LEFT JOIN books ON book_in_collection.book_id = books.book_id
		LEFT JOIN authors ON books.author_id = authors.author_id
    `

	// add to the wehre clause if it was present in the request args
    whereClauses := []string{}
    if c.CollectionID != nil {
        whereClauses = append(whereClauses, fmt.Sprintf("collections.collection_id = %d", *c.CollectionID))
    }
    if c.CollectionName != nil {
        whereClauses = append(whereClauses, fmt.Sprintf("collections.collection_name = '%s'", *c.CollectionName))
    }

    if len(whereClauses) > 0 {
        query += "WHERE " + strings.Join(whereClauses, " AND ")
    }

    query += "ORDER BY collections.collection_id"

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

	// add book structs to each collection struct, if that collection has a book related to it
    var currentCollection *Collection
    for rows.Next() {
        var book Book
        var collection Collection
		var bookID sql.NullInt64 // so we can scan even if there is no book associated to the collection
		var title sql.NullString // so we can scan even if there is no book associated to the collection
		var cDate sql.NullTime // so we can scan even if there is no book associated to the collection
		var author sql.NullString // so we can scan even if there is no book associated to the collection
        
        err := rows.Scan(&collection.CollectionID, &collection.CollectionName, &collection.CreationDate, &bookID, &title, &cDate, &author)
        if err != nil {
            return nil, err
        }

		// If this is a new collection, add it to the list
        if currentCollection == nil || currentCollection.CollectionID != collection.CollectionID {
            collections = append(collections, collection)
            currentCollection = &collections[len(collections)-1]
        }

		// create book object
		if bookID.Valid {
			book.BookID = int(bookID.Int64)
		}

		if title.Valid {
			book.Title = title.String
		}

		if author.Valid {
			book.Author = author.String
		}

		if cDate.Valid {
			book.CreationDate = cDate.Time
		}

        //If there is a related book, add it to the current collection
		if book.BookID != 0 {
			currentCollection.CollectionBooks = append(currentCollection.CollectionBooks, book)
		}
    }

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

func AddBookToCollection(db *sql.DB, a AddBookToCollectionArgs) (*Collection, *Book, error) {
	var collection *Collection
    var book *Book
    var err error

	// check if there is a book with the chosen ID
    if a.BookID != nil {
        books, err := ListBooks(db, BookArgs{BookID: a.BookID})
		if err != nil{
			return nil, nil, err
		}
        book = &books[0]
    } else {
        err = errors.New("choose the book to add to the collection and insert its ID number")
		return nil, nil, err
    }

	// check if there is a collection with the chosen ID
    if a.CollectionID != nil {
        collections, err := ListCollections(db, CollectionArgs{CollectionID: a.CollectionID})
		if err != nil{
			return nil, nil, err
		}
        collection = &collections[0]
    } else {
        err = errors.New("choose a collection to have the book added to its ID number")
		return nil, nil, err
    }

	err = db.QueryRow("INSERT INTO book_in_collection (book_id, collection_id) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING book_id, collection_id", book.BookID, collection.CollectionID).Scan(&book.BookID, &collection.CollectionID) // errors are deferred until Row's Scan method is called
    if err != nil {
        if err == sql.ErrNoRows{
            err = errors.New("book already in this collection")
        }
        return nil, nil, err
	}
	fmt.Printf("Book %s added to collection %s\n", book.Title, collection.CollectionName)

	return collection, book, nil
}
