package main

import "time"

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