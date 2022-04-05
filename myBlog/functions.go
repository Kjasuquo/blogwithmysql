package myBlog

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"html/template"
	"http/db_client"
	"net/http"
	"strings"
)

//Error handling
func Error(e error) {
	if e != nil {
		fmt.Println(e)
		return
	}
}

//Register initializes all my commands and html pages, and it is called in the main
func Register(router *chi.Mux) {
	router.Get("/", indexhandler)
	router.Get("/addpost", getContenthandler)
	router.Post("/addpost", postContenthandler)
	router.Get("/update/{Id}", updateByIdhandler)
	router.Post("/update/{Id}", postupdateByIdhandler)
	router.Get("/del/{Id}", deleteByIdhandler)
	router.Get("/readmore/{Id}", readMore)

}

//To get index page
func indexhandler(w http.ResponseWriter, r *http.Request) {

	//Scan whatever is in your DB into your DataStructure slice
	db_client.Scan()

	//This points to the html location
	t, e := template.ParseFiles("templat/index.html")
	Error(e)

	var f db_client.Data
	var y string
	var count int
	var g []db_client.Data
	for _, v := range db_client.DataStructure {
		text := strings.Split(v.Content, " ")
		for i := 0; i < len(text); i++ {
			count++
			if count == 100 {
				y += "..."
				break
			}
			y += text[i] + " "
		}
		f.Id = v.Id
		f.Title = v.Title
		f.Content = y
		f.Status = v.Status

		g = append(g, f)
		y = ""
		count = 0
	}

	//This writes whatever is in the DataStructure database to the html file
	e = t.Execute(w, g)
	if e != nil {
		fmt.Println(e)
	}

	db_client.DataStructure = nil
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
	f := db_client.Data{}

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

	db_client.InsertToDb(f)

	//redirect your page back to the index/home page when done (on a click)
	http.Redirect(w, r, "/", 302)
}

//To get the content on a page when Edit is clicked
func updateByIdhandler(w http.ResponseWriter, r *http.Request) {

	e := r.ParseForm()
	if e != nil {
		fmt.Println(e)
	}

	ID := chi.URLParam(r, "Id")

	//This points to the html location
	t, e := template.ParseFiles("templat/editpost.html")
	Error(e)

	//Calls or writes the item inside that database in the html file/template where it is called
	e = t.Execute(w, db_client.EditDb(ID))
	Error(e)
}

//After getting the content to be edited in html page, after editing, it will create a new item of it and store in database
func postupdateByIdhandler(w http.ResponseWriter, r *http.Request) {

	e := r.ParseForm()
	Error(e)
	id := chi.URLParam(r, "Id")

	//Gets the id/name of the form components
	title := r.PostForm.Get("tit")
	content := r.PostForm.Get("con")

	db_client.PostEditDb(title, content, id)

	//redirect your page back to the index/home page when done (on a click)
	http.Redirect(w, r, "/", 302)
}

//To delete each post
func deleteByIdhandler(w http.ResponseWriter, r *http.Request) {

	ID := chi.URLParam(r, "Id")

	db_client.DeletePost(ID)

	//redirect your page back to the index/home page when done (on a click)
	http.Redirect(w, r, "/", 302)
}

func readMore(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "Id")

	t, e := template.ParseFiles("templat/readmore.html")
	Error(e)

	//Calls or writes the item inside that database in the html file/template where it is called
	e = t.Execute(w, db_client.EditDb(ID))
	Error(e)

}
