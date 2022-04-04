package db_client

import (
	"database/sql"
	"fmt"
)

//var DBClient *sql.DB

func CreateAndOpen() {

	db, err := sql.Open("mysql", "root:flyn!GG@01@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS blogDB")
	if err != nil {
		panic(err)
	}
	db.Close()

	db, err = sql.Open("mysql", "root:flyn!GG@01@tcp(127.0.0.1:3306)/blogDB")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	var query string
	query = "CREATE TABLE IF NOT EXISTS Data(Id VARCHAR(500),Title VARCHAR(500),Content VARCHAR(1000),Status BOOLEAN DEFAULT TRUE )"
	if err != nil {
		panic(err)
	}
	create, err := db.Exec(query)

	//DBClient = db

	fmt.Println(create)

}
