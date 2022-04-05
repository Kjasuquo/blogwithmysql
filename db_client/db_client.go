package db_client

import (
	"database/sql"
	"fmt"
)

// DBClient that is accessible in other packages
var DBClient *sql.DB

// Data Structure containing everything I need to manipulate my data/content
type Data struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  bool   `json:"status"`
}

//DataStructure is a database in which my data are stored both before and after manipulation
var DataStructure []Data

// CreateAndOpen creates, Opens and keeps the DB open. Initialized in the main
func CreateAndOpen() {
	db, err := sql.Open("mysql", "root:flyn!GG@01@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS blogDB")
	if err != nil {
		panic(err)
	}

	//db.Close()

	db, err = sql.Open("mysql", "root:flyn!GG@01@tcp(127.0.0.1:3306)/blogDB")
	if err != nil {
		panic(err)
	}
	//
	//defer db.Close()

	var query string
	query = "CREATE TABLE IF NOT EXISTS Data(Id VARCHAR(500),Title VARCHAR(500),Content TEXT,Status BOOLEAN DEFAULT TRUE )"
	if err != nil {
		panic(err)
	}
	create, err := db.Exec(query)

	DBClient = db

	fmt.Println(create)
}

//Scan scans whatever in the DB into the DataStructure slice for easy manipulation.
func Scan() {
	st := "SELECT * FROM Data"
	rows, err := DBClient.Query(st)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var r Data
		err := rows.Scan(&r.Id, &r.Title, &r.Content, &r.Status)
		if err != nil {
			fmt.Println(err)
		}
		DataStructure = append(DataStructure, r)
	}
}

//InsertToDb inserts values from data into tHe DB
func InsertToDb(f Data) {
	insert, err := DBClient.Query("INSERT INTO Data (Id, Title, Content, STATUS) VALUES (?,?,?,?)", (f.Id), (f.Title), (f.Content), (f.Status))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(insert)
}

//EditDb helps to get an info/row with a particular ID form the DB
func EditDb(ID string) Data {
	row := DBClient.QueryRow("SELECT * FROM blogDB.Data WHERE Id = ?;", ID)

	var d Data

	er := row.Scan(&d.Id, &d.Title, &d.Content, &d.Status)
	if er != nil {
		fmt.Println(er)
	}
	return d
}

//PostEditDb helps to post whatever edited back into the DB
func PostEditDb(title string, content string, id string) {
	upst := "UPDATE `blogDB`.`Data` SET `Title` = ?, `Content` = ? WHERE (`Id`=?);"

	st, er := DBClient.Prepare(upst)
	if er != nil {
		fmt.Println(er)
	}

	defer st.Close()

	var res sql.Result
	res, er = st.Exec(title, content, id)
	rowAff, _ := res.RowsAffected()
	fmt.Println("rows affected:", rowAff)
}

//DeletePost helps in deleting from DB any row with a particular ID
func DeletePost(ID string) {
	del, err := DBClient.Prepare("DELETE FROM `blogDB`.`Data` WHERE (`Id` = ?);")
	if err != nil {
		fmt.Println(err)
	}
	defer del.Close()
	var res sql.Result
	res, err = del.Exec(ID)
	rowAff, _ := res.RowsAffected()
	fmt.Println("rows affected:", rowAff)
}
