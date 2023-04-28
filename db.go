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
