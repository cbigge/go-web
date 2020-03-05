package main

import (
	"html/template"
	"net/http"

	"github.com/cbigge/go-web/views"

	"github.com/gorilla/mux"
)

var homeView *views.View
var contactView *views.View

func main() {
	var err error
	homeView, err = views.NewView("views/home.gohtml")
	if err != nil {
		panic(err)
	}
	contactView, err = views.NewView("views/contact.gohtml")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := homeView.Template.ExecuteTemplate(w, homeView.Layout, nil)
	if err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := contactView.Template.ExecuteTemplate(w, homeView.Layout, nil)
	if err != nil {
		panic(err)
	}
}
