package main

import (
	"html/template"
	"net/http"
	"time"
)

type Article struct {
	Path         string
	Title        string
	Content      string
	CreationDate time.Time
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/base.html", "./templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/home", homeHandler)
	http.ListenAndServe(":8080", nil)
}
