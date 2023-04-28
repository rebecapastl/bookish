package main

import (
	"flag"
	"time"
)

type Config struct {
	Database struct {
		URL string `yaml:"url"`
	} `yaml:"database"`
}

type AuthorArgs struct {
	Name *string `json:"name"`
}

type Author struct {
	AuthorID   int       `json:"author_id"`
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creation_date"`
}

type BookArgs struct {
	BookID  *int    `json:"book_id"`
	Title    *string `json:"title"`
	Author *string `json:"author"`
}

type Book struct {
	BookID      int       `json:"book_id"`
	Title        string    `json:"title"`
	Author     string    `json:"author"`
	CreationDate time.Time `json:"creation_date"`
}

type CollectionArgs struct {
	CollectionID	*int	`json:"collection_id"`
	CollectionName *string `json:"collection_name"`
}

type Collection struct {
	CollectionID   int	`json:"collection_id"`
	CollectionName string `json:"collection_name"`
	CreationDate	time.Time `json:"creation_date"`
	CollectionBooks	[]Book
}

type AddBookToCollectionArgs struct {
	BookID      *int `json:"book_id"`
	CollectionID *int `json:"collection_id"`
}

// Command represents a command with its subcommands and associated flags.
type Command struct {
	name        string
	description string
	subcommands []*Subcommand
}

// Subcommand represents a subcommand with its associated flags.
type Subcommand struct {
	name        string
	description string
	flags       *flag.FlagSet
}
