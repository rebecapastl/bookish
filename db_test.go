package main_test

import (
	"database/sql"
	"testing"
	"time"

	"bookish"

	"github.com/stretchr/testify/suite"
)

type DbTestSuite struct {
    suite.Suite
    db *sql.DB
}

func (suite *DbTestSuite) SetupTest() {
    // Connect to the test db
	testConfig, err := main.LoadTestConfig()
	if err != nil {
        suite.T().Fatal(err)
    }
	db, err := sql.Open("postgres", testConfig.Database.URL)
    if err != nil {
        suite.T().Fatal(err)
    }
    suite.db = db

	main.CreateTables(suite.db)
}

func (suite *DbTestSuite) TearDownTest() {
    _, err := suite.db.Exec("DROP TABLE IF EXISTS book_in_collection")
    if err != nil {
        suite.T().Fatal(err)
    }
    
    _, err = suite.db.Exec("DROP TABLE IF EXISTS collections")
    if err != nil {
        suite.T().Fatal(err)
    }

    _, err = suite.db.Exec("DROP TABLE IF EXISTS books")
    if err != nil {
        suite.T().Fatal(err)
    }

    _, err = suite.db.Exec("DROP TABLE IF EXISTS authors")
    if err != nil {
        suite.T().Fatal(err)
    }

	suite.db.Close()
}


func (suite *DbTestSuite) TestCreateAuthor() {
	// Setup
	authorName := "J. R. R. Tolkien"

	// Function to test
	author, err := main.CreateAuthor(suite.db, main.AuthorArgs{Name: &authorName})

	// Verification
	suite.NoError(err)
	suite.NotNil(author)
	suite.Equal(authorName, author.Name)
}

func (suite *DbTestSuite) TestCreateAuthor_DuplicateAuthor() {
	// Setup
	authorName := "J. R. R. Tolkien"
	_, err := main.CreateAuthor(suite.db, main.AuthorArgs{Name: &authorName})
	suite.NoError(err)

	// Function to test
	duplicateAuthor, err := main.CreateAuthor(suite.db, main.AuthorArgs{Name: &authorName})

	// Verification
	suite.Error(err)
	suite.Equal("author already exists in the database", err.Error())
	suite.Nil(duplicateAuthor)
}

func (suite *DbTestSuite) TestListAuthors() {
    // Setup
    _, err := suite.db.Exec("INSERT INTO authors (name, creation_date) VALUES ($1, $2)", "J. R. R. Tolkien", time.Now().UTC())
    suite.NoError(err)
    _, err = suite.db.Exec("INSERT INTO authors (name, creation_date) VALUES ($1, $2)", "Octavia E. Butler", time.Now().UTC())
    suite.NoError(err)

	expectedAuthors := []main.Author{
		{AuthorID: 1, Name: "J. R. R. Tolkien", CreationDate: time.Now().UTC()},
		{AuthorID: 2, Name: "Octavia E. Butler", CreationDate: time.Now().UTC()},
	}

	// Function to test
    authors, err := main.ListAuthors(suite.db, main.AuthorArgs{})
    
	// Verification
	suite.NoError(err)
	suite.Len(authors, 2)
    for index, author := range authors {
		suite.Equal(expectedAuthors[index].AuthorID, author.AuthorID)
		suite.Equal(expectedAuthors[index].Name, author.Name)
		suite.Equal(expectedAuthors[index].CreationDate.Format("2006-01-02"), author.CreationDate.Format("2006-01-02"))

	}

}

func (suite *DbTestSuite) TestListAuthors_ByName() {
    // Setup
    _, err := suite.db.Exec("INSERT INTO authors (name, creation_date) VALUES ($1, $2)", "J. R. R. Tolkien", time.Now().UTC())
    suite.NoError(err)
    _, err = suite.db.Exec("INSERT INTO authors (name, creation_date) VALUES ($1, $2)", "Octavia E. Butler", time.Now().UTC())
    suite.NoError(err)

    author := "J. R. R. Tolkien"

	// Function to test
    authors, err := main.ListAuthors(suite.db, main.AuthorArgs{Name: &author})
    
    // Verification
	suite.NoError(err)
    suite.Len(authors, 1)
    suite.Equal(author, authors[0].Name)
}

func (suite *DbTestSuite) TestListAuthors_NoAuthor() {
    // Setup
    _, err := suite.db.Exec("INSERT INTO authors (name, creation_date) VALUES ($1, $2)", "J. R. R. Tolkien", time.Now().UTC())
    suite.NoError(err)
    _, err = suite.db.Exec("INSERT INTO authors (name, creation_date) VALUES ($1, $2)", "Octavia E. Butler", time.Now().UTC())
    suite.NoError(err)

    author := "Juliet Marillier"

	// Function to test
    authors, err := main.ListAuthors(suite.db, main.AuthorArgs{Name: &author})

    // Verification
	suite.NoError(err)
	suite.Empty(authors)
}

func (suite *DbTestSuite) TestCreateBook_WithAuthor() {
    // Setup
	author := "J. R. R. Tolkien"
	bookName := "Book 1"

	// Function to test
    book, err := main.CreateBook(suite.db, main.BookArgs{Title:&bookName, Author: &author})

    // Verification
	suite.NoError(err)
	suite.Equal(bookName, book.Title)
    suite.Equal(author, book.Author)
}

func (suite *DbTestSuite) TestCreateBook_WithoutAuthor() {
    // Setup
	bookName := "Book 1"

	// Function to test
    book, err := main.CreateBook(suite.db, main.BookArgs{Title:&bookName})

    // Verification
	suite.NoError(err)
	suite.Equal(bookName, book.Title)
    suite.Equal("anonymous", book.Author)
}

func (suite *DbTestSuite) TestCreateBook_DuplicateBook() {
	// Setup
	bookName := "Book 1"
    _, err := main.CreateBook(suite.db, main.BookArgs{Title:&bookName})
	suite.NoError(err)

	// Function to test
	duplicateBook, err := main.CreateBook(suite.db, main.BookArgs{Title:&bookName})

	// Verification
	suite.Error(err)
	suite.Equal("book already exists in the database", err.Error())
	suite.Nil(duplicateBook)
}


func (suite *DbTestSuite) TestListBooks(){
    // Setup
    _, err := suite.db.Exec("INSERT INTO authors (name) VALUES ('J. R. R. Tolkien'), ('Octavia E. Butler'), ('G. R. R. Martin')")
    suite.NoError(err)
    _, err = suite.db.Exec("INSERT INTO books (title, author_id) VALUES ('Book 1', 1), ('Book 2', 2), ('Book 3', 3)")
    suite.NoError(err)

	expectedBooks := []main.Book{
		{BookID: 1, Title: "Book 1", Author: "J. R. R. Tolkien", CreationDate: time.Now().UTC()},
		{BookID: 2, Title: "Book 2", Author: "Octavia E. Butler", CreationDate: time.Now().UTC()},
		{BookID: 3, Title: "Book 3", Author: "G. R. R. Martin", CreationDate: time.Now().UTC()},
	}

	// Function to test
	books, err := main.ListBooks(suite.db, main.BookArgs{})
	suite.NoError(err)

	// Verification
	suite.Len(books, 3)
    for index, book := range books {
		suite.Equal(expectedBooks[index].BookID, book.BookID)
		suite.Equal(expectedBooks[index].Author, book.Author)
		suite.Equal(expectedBooks[index].CreationDate.Format("2006-01-02"), book.CreationDate.Format("2006-01-02"))

	}

}

func (suite *DbTestSuite) TestListBooks_ByName(){
    // Setup
    _, err := suite.db.Exec("INSERT INTO authors (name) VALUES ('J. R. R. Tolkien'), ('Octavia E. Butler'), ('G. R. R. Martin')")
    suite.NoError(err)
    _, err = suite.db.Exec("INSERT INTO books (title, author_id) VALUES ('Book 1', 1), ('Book 1', 2), ('Book 3', 3)")
    suite.NoError(err)

	expectedBooks := []main.Book{
		{BookID: 1, Title: "Book 1", Author: "J. R. R. Tolkien", CreationDate: time.Now().UTC()},
		{BookID: 2, Title: "Book 1", Author: "Octavia E. Butler", CreationDate: time.Now().UTC()},
	}

	// Function to test
	bookName := "Book 1"
	books, err := main.ListBooks(suite.db, main.BookArgs{Title: &bookName})

	// Verification
	suite.NoError(err)
	suite.Len(books, 2)
    for index, book := range books {
		suite.Equal(expectedBooks[index].BookID, book.BookID)
		suite.Equal(expectedBooks[index].Author, book.Author)
		suite.Equal(expectedBooks[index].CreationDate.Format("2006-01-02"), book.CreationDate.Format("2006-01-02"))

	}

}

func (suite *DbTestSuite) TestListBooks_ByAuthor(){
    // Setup
    _, err := suite.db.Exec("INSERT INTO authors (name) VALUES ('J. R. R. Tolkien'), ('Octavia E. Butler'), ('G. R. R. Martin')")
    suite.NoError(err)
    _, err = suite.db.Exec("INSERT INTO books (title, author_id) VALUES ('Book 1', 1), ('Book 2', 3), ('Book 3', 3)")
    suite.NoError(err)

	expectedBooks := []main.Book{
		{BookID: 2, Title: "Book 2", Author: "G. R. R. Martin", CreationDate: time.Now().UTC()},
		{BookID: 3, Title: "Book 3", Author: "G. R. R. Martin", CreationDate: time.Now().UTC()},
	}

	// Function to test
	authorName := "G. R. R. Martin"
	books, err := main.ListBooks(suite.db, main.BookArgs{Author: &authorName})

	// Verification
	suite.NoError(err)
	suite.Len(books, 2)
    for index, book := range books {
		suite.Equal(expectedBooks[index].BookID, book.BookID)
		suite.Equal(expectedBooks[index].Author, book.Author)
		suite.Equal(expectedBooks[index].CreationDate.Format("2006-01-02"), book.CreationDate.Format("2006-01-02"))

	}

}

func (suite *DbTestSuite) TestListBooks_ById(){
    // Setup
    _, err := suite.db.Exec("INSERT INTO authors (name) VALUES ('J. R. R. Tolkien'), ('Octavia E. Butler'), ('G. R. R. Martin')")
    suite.NoError(err)
    _, err = suite.db.Exec("INSERT INTO books (title, author_id) VALUES ('Book 1', 1), ('Book 2', 2), ('Book 3', 3)")
    suite.NoError(err)

	expectedBook := main.Book{BookID: 1, Title: "Book 1", Author: "J. R. R. Tolkien", CreationDate: time.Now().UTC()}
	
	// Function to test
	id := 1
	book, err := main.ListBooks(suite.db, main.BookArgs{BookID: &id})
	
	// Verification
	suite.NoError(err)
	suite.Len(book, 1)
	suite.Equal(expectedBook.BookID, book[0].BookID)
	suite.Equal(expectedBook.Author, book[0].Author)
	suite.Equal(expectedBook.CreationDate.Format("2006-01-02"), book[0].CreationDate.Format("2006-01-02"))
}

func (suite *DbTestSuite) TestListBooks_NoBook() {
    // Setup
    _, err := suite.db.Exec("INSERT INTO authors (name) VALUES ('J. R. R. Tolkien'), ('Octavia E. Butler'), ('G. R. R. Martin')")
    suite.NoError(err)
    _, err = suite.db.Exec("INSERT INTO books (title, author_id) VALUES ('Book 1', 1), ('Book 2', 2), ('Book 3', 3)")
    suite.NoError(err)

    bookName := "Book 4"

	// Function to test
    books, err := main.ListBooks(suite.db, main.BookArgs{Title: &bookName})

    // Verification
	suite.Error(err)
	suite.Equal("no books with the chosen specification", err.Error())
	suite.Empty(books)
}

func (suite *DbTestSuite) TestCreateCollection() {
	// Setup
	collectionName := "My Collection"

	// Function to test
	collection, err := main.CreateCollection(suite.db, main.CollectionArgs{CollectionName: &collectionName})

	// Verification
	suite.NoError(err)
	suite.NotNil(collection)
	suite.Equal(collectionName, collection.CollectionName)
}

func (suite *DbTestSuite) TestCreateCollection_DuplicateCollection() {
	// Setup
	collectionName := "My Collection"
	_, err := main.CreateCollection(suite.db, main.CollectionArgs{CollectionName: &collectionName})
	suite.NoError(err)

	// Function to test
	duplicateCollection, err := main.CreateCollection(suite.db, main.CollectionArgs{CollectionName: &collectionName})

	// Verification
	suite.Error(err)
	suite.Equal("collection already exists in the database", err.Error())
	suite.Nil(duplicateCollection)
}

func (suite *DbTestSuite) TestListCollection() {
    // Setup
    _, err := suite.db.Exec("INSERT INTO collections (collection_name, creation_date) VALUES ($1, $2)", "My Collection 1", time.Now().UTC())
    suite.NoError(err)
    _, err = suite.db.Exec("INSERT INTO collections (collection_name, creation_date) VALUES ($1, $2)", "My Collection 2", time.Now().UTC())
    suite.NoError(err)

	expectedCollections := []main.Collection{
		{CollectionID: 1, CollectionName: "My Collection 1", CreationDate: time.Now().UTC(), CollectionBooks: nil},
		{CollectionID: 2, CollectionName: "My Collection 2", CreationDate: time.Now().UTC(), CollectionBooks: nil},
	}

	// Function to test
    collections, err := main.ListCollections(suite.db, main.CollectionArgs{})
    
    // Verification
	suite.NoError(err)
    suite.Len(collections, 2)
	for index, collection := range collections {
		suite.Equal(expectedCollections[index].CollectionID, collection.CollectionID)
		suite.Equal(expectedCollections[index].CollectionName, collection.CollectionName)
		suite.Equal(expectedCollections[index].CreationDate.Format("2006-01-02"), collection.CreationDate.Format("2006-01-02"))
		suite.Equal(expectedCollections[index].CollectionBooks, collection.CollectionBooks)

	}
}

func (suite *DbTestSuite) TestListCollection_ByName() {
    // Setup
    _, err := suite.db.Exec("INSERT INTO collections (collection_name, creation_date) VALUES ($1, $2)", "My Collection 1", time.Now().UTC())
    suite.NoError(err)
    _, err = suite.db.Exec("INSERT INTO collections (collection_name, creation_date) VALUES ($1, $2)", "My Collection 2", time.Now().UTC())
    suite.NoError(err)

    collectionName := "My Collection 2"

	// Function to test
    collections, err := main.ListCollections(suite.db, main.CollectionArgs{CollectionName: &collectionName})
    
    // Verification
	suite.NoError(err)
    suite.Len(collections, 1)
    suite.Equal(collectionName, collections[0].CollectionName)
}

func (suite *DbTestSuite) TestListCollection_NoCollection() {
    // Setup
    _, err := suite.db.Exec("INSERT INTO collections (collection_name, creation_date) VALUES ($1, $2)", "My Collection 1", time.Now().UTC())
    suite.NoError(err)
    _, err = suite.db.Exec("INSERT INTO collections (collection_name, creation_date) VALUES ($1, $2)", "My Collection 2", time.Now().UTC())
    suite.NoError(err)

    collectionName := "My Collection 3"

	// Function to test
    collections, err := main.ListCollections(suite.db, main.CollectionArgs{CollectionName: &collectionName})

    // Verification
	suite.Error(err)
	suite.Equal("no collections with the chosen specification", err.Error())
	suite.Empty(collections)
}

func (suite *DbTestSuite) TestAddBookToCollection(){
	// Setup
	_, err := suite.db.Exec("INSERT INTO collections (collection_name, creation_date) VALUES ($1, $2)", "My Collection 1", time.Now().UTC())
    suite.NoError(err)
	_, err = suite.db.Exec("INSERT INTO authors (name) VALUES ('Octavia E. Butler')")
    suite.NoError(err)
	_, err = suite.db.Exec("INSERT INTO books (title, author_id) VALUES ('Book 1', 1)")
    suite.NoError(err)

	bookId := 1
	collectionId := 1
	addArgs := main.AddBookToCollectionArgs{BookID: &bookId, CollectionID: &collectionId}

	// Function to test
	addedCollection, addedBook, err := main.AddBookToCollection(suite.db, addArgs)

	// Verification
	suite.NoError(err)

	suite.Equal(1, addedBook.BookID)
	suite.Equal("Book 1", addedBook.Title)
	suite.Equal("Octavia E. Butler", addedBook.Author)

	suite.Equal(1, addedCollection.CollectionID)
	suite.Equal("My Collection 1", addedCollection.CollectionName)

	// Verify updated collection
	collections, err := main.ListCollections(suite.db, main.CollectionArgs{CollectionName: &addedCollection.CollectionName})
	suite.NoError(err)

	suite.Equal(1, collections[0].CollectionBooks[0].BookID)
	suite.Equal("Book 1", collections[0].CollectionBooks[0].Title)
	suite.Equal("Octavia E. Butler", collections[0].CollectionBooks[0].Author)
}

func (suite *DbTestSuite) TestAddBookToCollection_NoBook(){
	// Setup
	_, err := suite.db.Exec("INSERT INTO collections (collection_name, creation_date) VALUES ($1, $2)", "My Collection 1", time.Now().UTC())
    suite.NoError(err)

	collectionId := 1
	addArgs := main.AddBookToCollectionArgs{CollectionID: &collectionId}

	// Function to test
	addedCollection, addedBook, err := main.AddBookToCollection(suite.db, addArgs)

	// Verification
	suite.Error(err)
	suite.Nil(addedBook)
	suite.Nil(addedCollection)
	suite.Equal("choose the book to add to the collection and insert its ID number", err.Error())
}

func (suite *DbTestSuite) TestAddBookToCollection_NoCollection(){
	// Setup
	_, err := suite.db.Exec("INSERT INTO authors (name) VALUES ('Octavia E. Butler')")
    suite.NoError(err)
	_, err = suite.db.Exec("INSERT INTO books (title, author_id) VALUES ('Book 1', 1)")
    suite.NoError(err)

	bookId := 1
	addArgs := main.AddBookToCollectionArgs{BookID: &bookId}

	// Function to test
	addedCollection, addedBook, err := main.AddBookToCollection(suite.db, addArgs)

	// Verification
	suite.Error(err)
	suite.Nil(addedBook)
	suite.Nil(addedCollection)
	suite.Equal("choose a collection to have the book added to its ID number", err.Error())
}

func (suite *DbTestSuite) TestAddBookToCollection_NoBookExistsWithChosenID(){
	// Setup
	_, err := suite.db.Exec("INSERT INTO collections (collection_name, creation_date) VALUES ($1, $2)", "My Collection 1", time.Now().UTC())
    suite.NoError(err)

	bookId := 1
	collectionId := 1
	addArgs := main.AddBookToCollectionArgs{BookID: &bookId, CollectionID: &collectionId}

	// Function to test
	addedCollection, addedBook, err := main.AddBookToCollection(suite.db, addArgs)

	// Verification
	suite.Error(err)
	suite.Nil(addedBook)
	suite.Nil(addedCollection)
	suite.Equal("no books with the chosen specification", err.Error())
}

func (suite *DbTestSuite) TestAddBookToCollection_NoCollectionExistsWithChosenID(){
	// Setup
	_, err := suite.db.Exec("INSERT INTO authors (name) VALUES ('Octavia E. Butler')")
    suite.NoError(err)
	_, err = suite.db.Exec("INSERT INTO books (title, author_id) VALUES ('Book 1', 1)")
    suite.NoError(err)

	bookId := 1
	collectionId := 1
	addArgs := main.AddBookToCollectionArgs{BookID: &bookId, CollectionID: &collectionId}

	// Function to test
	addedCollection, addedBook, err := main.AddBookToCollection(suite.db, addArgs)

	// Verification
	suite.Error(err)
	suite.Nil(addedBook)
	suite.Nil(addedCollection)
	suite.Equal("no collections with the chosen specification", err.Error())
}



func TestDbTestSuite(t *testing.T) {
    suite.Run(t, new(DbTestSuite))
}