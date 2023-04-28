package main

import (
	"database/sql"
	"fmt"
	"os"

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
}