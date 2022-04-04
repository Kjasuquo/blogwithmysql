package myBlog

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"html/template"
	"net/http"
)

//Error handling
func Error(e error) {
	if e != nil {
		return
	}
}

// Data Structure containing everything I need to manipulate my data/content
type Data struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  bool   `json:"status"`
}

//DataStructure is a database in which my data are stored both before and after manipulation
var DataStructure []Data

//Register initializes all my commands and html pages and it is called in the main
func Register(router *chi.Mux) {
	router.Get("/", indexhandler)
	router.Get("/addpost", getContenthandler)
	router.Post("/addpost", postContenthandler)
	router.Get("/update/{Id}", updateByIdhandler)
	router.Post("/update/{Id}", postupdateByIdhandler)
	router.Get("/del/{Id}", deleteByIdhandler)

}

//To get index page
func indexhandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "root:flyn!GG@01@tcp(127.0.0.1:3306)/blogDB")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	st := "SELECT * FROM Data"
	rows, err := db.Query(st)
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

	//This points to the html location
	t, e := template.ParseFiles("templat/index.html")
	Error(e)

	//This writes whatever is in the DataStructure database to the html file
	e = t.Execute(w, DataStructure)
	Error(e)

	DataStructure = nil
}

//To get content page
func getContenthandler(w http.ResponseWriter, r *http.Request) {
	//This points to the html location

	t, e := template.ParseFiles("templat/addpost.html")
	Error(e)

	//This displays whatever on the html page
	e = t.Execute(w, nil)
	Error(e)
}

//To manipulate content page
func postContenthandler(w http.ResponseWriter, r *http.Request) {
	//creating an instance of a data struct
	f := Data{}

	//This gets/populates the content of the form
	e := r.ParseForm()
	Error(e)

	//Gets the id/name of the form components
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	//Filling the data struct instance
	f.Title = title
	f.Content = content
	f.Status = true
	f.Id = uuid.NewString() //new id is being populated for an item using google/uuid

	//Attach the f to the Data base so that it can be populated on the index page
	//DataStructure = append(DataStructure, f)

	db, err := sql.Open("mysql", "root:flyn!GG@01@tcp(127.0.0.1:3306)/blogDB")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	insert, err := db.Query("INSERT INTO Data (Id, Title, Content, STATUS) VALUES (?,?,?,?)", (f.Id), (f.Title), (f.Content), (f.Status))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(insert)

	//redirect your page back to the index/home page when done (on a click)
	http.Redirect(w, r, "/", 302)
}

//To get the content on a page when Edit is clicked
func updateByIdhandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "root:flyn!GG@01@tcp(127.0.0.1:3306)/blogDB")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	r.ParseForm()

	ID := chi.URLParam(r, "Id")

	row := db.QueryRow("SELECT * FROM blogDB.Data WHERE Id = ?;", ID)

	var d Data

	er := row.Scan(&d.Id, &d.Title, &d.Content, &d.Status)
	if err != nil {
		fmt.Println(er)
		return
	}

	//This points to the html location
	t, e := template.ParseFiles("templat/editpost.html")
	Error(e)

	//Calls or writes the item inside that database in the html file/template where it is called
	e = t.Execute(w, d)
	Error(e)
}

//After getting the content to be edited in html page, after editing, it will create a new item of it and store in database
func postupdateByIdhandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:flyn!GG@01@tcp(127.0.0.1:3306)/blogDB")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	e := r.ParseForm()
	Error(e)
	id := chi.URLParam(r, "Id")

	//Gets the id/name of the form components
	title := r.PostForm.Get("tit")
	content := r.PostForm.Get("con")

	upst := "UPDATE `blogDB`.`Data` SET `Title` = ?, `Content` = ? WHERE (`Id`=?);"

	st, e := db.Prepare(upst)
	if e != nil {
		fmt.Println(e)
	}

	defer st.Close()

	var res sql.Result
	res, err = st.Exec(title, content, id)
	rowAff, _ := res.RowsAffected()
	fmt.Println("rows affected:", rowAff)

	//redirect your page back to the index/home page when done (on a click)
	http.Redirect(w, r, "/", 302)
}

//To delete each post
func deleteByIdhandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "root:flyn!GG@01@tcp(127.0.0.1:3306)/blogDB")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//Whenever it is clicked, get the id of its element
	ID := chi.URLParam(r, "Id")

	del, err := db.Prepare("DELETE FROM `blogDB`.`Data` WHERE (`Id` = ?);")
	if err != nil {
		fmt.Println(err)
	}
	defer del.Close()
	var res sql.Result
	res, err = del.Exec(ID)
	rowAff, _ := res.RowsAffected()
	fmt.Println("rows affected:", rowAff)

	//redirect your page back to the index/home page when done (on a click)
	http.Redirect(w, r, "/", 302)
}
