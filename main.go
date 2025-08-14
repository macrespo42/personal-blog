package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

func getArticleById(id int) (Article, error) {
	fileName := fmt.Sprintf("%v.json", id)
	fileData, err := os.ReadFile(filepath.Join("./articles", fileName))
	if err != nil {
		return Article{}, err
	}

	var article Article
	err = json.Unmarshal(fileData, &article)
	if err != nil {
		return Article{}, err
	}
	return article, nil
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

func articleHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := template.ParseFiles("./templates/base.html", "./templates/article.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	idNb, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	article, err := getArticleById(idNb)
	err = t.Execute(w, article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/home", homeHandler)
	router.HandleFunc("/article/{id}", articleHandler)
	http.ListenAndServe(":8080", router)
}
