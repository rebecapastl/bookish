package main

import (
	"flag"
	"fmt"
	"os"
)

func createBookCommands() Command {
	var createBookTitle string
	var createBookAuthor string
	var listBookTitle string
	var listBookAuthor string
	var listBookId string

	// Define command-line interface
	bookCmd := Command{
		name:        "book",
		description: "Manage books in the database",
		subcommands: []*Subcommand{
			{
				name:        "create",
				description: "Create a new book",
				flags:       flag.NewFlagSet("create", flag.ExitOnError),
			},
			{
				name:        "list",
				description: "List all books",
				flags:       flag.NewFlagSet("list", flag.ExitOnError),
			},
		},
	}

	// Define flags for the 'create' subcommand of the 'book' command
	createBookCmd := bookCmd.subcommands[0].flags
	createBookCmd.StringVar(&createBookTitle, "t", "", "Title of the book")
	createBookCmd.StringVar(&createBookAuthor, "a", "", "Name of the author")

	// Define flags for the 'list' subcommand of the 'book' command
	listBookCmd := bookCmd.subcommands[1].flags
	listBookCmd.StringVar(&listBookTitle, "t", "", "Title of the book")
	listBookCmd.StringVar(&listBookAuthor, "a", "", "Name of the author")
	listBookCmd.StringVar(&listBookId, "i", "", "Id of the book")

	return bookCmd
}

func createCollectionCommands() Command {
	var createCollectionName string
	var listCollectionName string
	var listCollectionId string

	var addCollectionId string
	var addBookId string

	collectionCmd := Command{
		name:        "collection",
		description: "Manage collections in the database",
		subcommands: []*Subcommand{
			{
				name:        "create",
				description: "Create a new collection",
				flags:       flag.NewFlagSet("create", flag.ExitOnError),
			},
			{
				name:        "list",
				description: "List all collections",
				flags:       flag.NewFlagSet("list", flag.ExitOnError),
			},
			{
				name:        "add",
				description: "Add a book to a collection",
				flags:       flag.NewFlagSet("add", flag.ExitOnError),
			},
		},
	}

	// Define flags for the 'create' subcommand of the 'collection' command
	createCollectionCmd := collectionCmd.subcommands[0].flags
	createCollectionCmd.StringVar(&createCollectionName, "n", "", "Name of the collection")

	// Define flags for the 'list' subcommand of the 'collection' command
	listCollectionCmd := collectionCmd.subcommands[1].flags
	listCollectionCmd.StringVar(&listCollectionName, "n", "", "Name of the collection")
	listCollectionCmd.StringVar(&listCollectionId, "i", "", "Id of the collection")

	// Define flags for the 'list' subcommand of the 'collection' command
	addCollectionCmd := collectionCmd.subcommands[2].flags
	addCollectionCmd.StringVar(&addCollectionId, "i", "", "Id of the collection")
	addCollectionCmd.StringVar(&addBookId, "bi", "", "Id of the book to be added")

	return collectionCmd
}

func CLIcommands() {

	bookCmd := createBookCommands()
	createBookCmd := bookCmd.subcommands[0].flags
	listBookCmd := bookCmd.subcommands[1].flags

	collectionCmd := createCollectionCommands()
	createCollectionCmd := collectionCmd.subcommands[0].flags
	listCollectionCmd := collectionCmd.subcommands[1].flags
	addCollectionCmd := collectionCmd.subcommands[2].flags

	// Parse command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: books-database <command> [<args>]")
		fmt.Println("Commands:")
		fmt.Println("\tbook create\tCreate a new book")
		fmt.Println("\tbook list\t\t\tList all books")
		fmt.Println("\tcollection create\t\tCreate a new collection")
		fmt.Println("\tcollection list\t\tList all collections")
		fmt.Println("\tcollection add\t\tAdd a book to a collection")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "book":
		// Parse subcommand arguments
		if len(os.Args) < 3 {
			fmt.Println("No book title set, book not created")
			fmt.Println("Usage: books-database book <subcommand> [<args>]")
			fmt.Println("Subcommands:")
			fmt.Println("\tcreate\tCreate a new book")
			fmt.Println("\tlist\t\t\tList all books")
			os.Exit(1)
		}

		var bookArgs BookArgs
		switch os.Args[2] {
		case "create":
			bookCmd.subcommands[0].flags.Parse(os.Args[3:])
			var tFlag *string 
			var aFlag *string 
			
			if createBookCmd.Lookup("t").Value.String() != "" {
				tFlagString := createBookCmd.Lookup("t").Value.String() //not addressable
				tFlag = &tFlagString
			}
			if createBookCmd.Lookup("a").Value.String() != "" {
				aFlagString := createBookCmd.Lookup("a").Value.String() //not addressable
				aFlag = &aFlagString
				
			}

			bookArgs = BookArgs{Title: tFlag, Author: aFlag}

			result, err := CreateBook(db, bookArgs)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Creating book with title %s\n", result.Title)
			}

		case "list":
			bookCmd.subcommands[1].flags.Parse(os.Args[3:])
			var tFlag *string 
			var aFlag *string 
			var iFlag *int 

			if listBookCmd.Lookup("t").Value.String() != "" {
				tFlagString := listBookCmd.Lookup("t").Value.String() //not addressable
				tFlag = &tFlagString
			}
			if listBookCmd.Lookup("a").Value.String() != "" {
				aFlagString := listBookCmd.Lookup("a").Value.String() //not addressable
				aFlag = &aFlagString
			}
			if listBookCmd.Lookup("i").Value.String() != "" {
				iString := listBookCmd.Lookup("i").Value.String()
				iFlag, err = SanitizeIdNumber(&iString) //not addressable
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			
			bookArgs = BookArgs{BookID: iFlag, Title: tFlag, Author: aFlag}
			result, err := ListBooks(db, bookArgs)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
		default:
			fmt.Println("Invalid subcommand. Expected 'create' or 'list'.")
			os.Exit(1)
		}

	case "collection":
		// Parse subcommand arguments
		if len(os.Args) < 3 {
			fmt.Println("No collenction name set, collection not created")
			fmt.Println("Usage: books-database collection <subcommand> [<args>]")
			fmt.Println("Subcommands:")
			fmt.Println("\tcreate\tCreate a new author")
			fmt.Println("\tlist\tList all authors")
			fmt.Println("\add\tAdd a book to a collection")
			os.Exit(1)
		}

		var collectionArgs CollectionArgs
		switch os.Args[2] {
		case "create":
			collectionCmd.subcommands[0].flags.Parse(os.Args[3:])
			var nFlag *string 
			
			if createCollectionCmd.Lookup("n").Value.String() != "" {
				nFlagString := createCollectionCmd.Lookup("n").Value.String() //not addressable
				nFlag = &nFlagString
			}

			collectionArgs = CollectionArgs{ CollectionName: nFlag}

			result, err := CreateCollection(db, collectionArgs)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Creating collection with name %s\n", result.CollectionName)
			}

		case "list":
			collectionCmd.subcommands[1].flags.Parse(os.Args[3:])
			var nFlag *string 
			var iFlag *int 

			if listCollectionCmd.Lookup("n").Value.String() != "" {
				nFlagString := listCollectionCmd.Lookup("n").Value.String() //not addressable
				nFlag = &nFlagString
			}
			if listCollectionCmd.Lookup("i").Value.String() != "" {
				iString := listCollectionCmd.Lookup("i").Value.String()
				iFlag, err = SanitizeIdNumber(&iString) //not addressable
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			
			collectionArgs = CollectionArgs{CollectionID: iFlag,CollectionName: nFlag}
			result, err := ListCollections(db, collectionArgs)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}

		case "add":
			collectionCmd.subcommands[2].flags.Parse(os.Args[3:])
			var iFlag *int 
			var biFlag *int 
			
			if addCollectionCmd.Lookup("i").Value.String() != "" {
				iFlagString := addCollectionCmd.Lookup("i").Value.String() //not addressable
				iFlag, err = SanitizeIdNumber(&iFlagString) //not addressable
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			if addCollectionCmd.Lookup("bi").Value.String() != "" {
				biFlagString := addCollectionCmd.Lookup("bi").Value.String() //not addressable
				biFlag, err = SanitizeIdNumber(&biFlagString) //not addressable
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			addArgs := AddBookToCollectionArgs{ BookID: biFlag, CollectionID: iFlag}

			_, _, err = AddBookToCollection(db, addArgs)
			if err != nil {
				fmt.Println(err)
			}

		default:
			fmt.Println("Invalid subcommand. Expected 'create' or 'list'.")
			os.Exit(1)
		}

	default:
		fmt.Println("Invalid command. Expected 'book' or 'colletion'.")
		os.Exit(1)
	}

}