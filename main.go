package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Article struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CreationDate time.Time `json:"creationDate"`
}

func getArticles() ([]Article, error) {
	files, err := os.ReadDir("./articles")
	if err != nil {
		return nil, err
	}

	var articles []Article

	for _, file := range files {
		fileData, err := os.ReadFile(filepath.Join("./articles", file.Name()))
		if err != nil {
			return nil, err
		}
		var article Article
		err = json.Unmarshal(fileData, &article)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/base.html", "./templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	articles, err := getArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/home", homeHandler)
	http.ListenAndServe(":8080", router)
}
