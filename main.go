package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Article struct {
	Id           string    `json:"id"`
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

func getArticleById(id string) (Article, error) {
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
	article, err := getArticleById(id)
	err = t.Execute(w, article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func newArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	id := uuid.NewString()
	creationDate := time.Now()
	title := r.FormValue("title")
	content := r.FormValue("content")

	if title != "" && content != "" {
		article := Article{
			Id:           id,
			CreationDate: creationDate,
			Title:        title,
			Content:      content,
		}

		fileName := fmt.Sprintf("%v.json", id)
		articlePath := filepath.Join("./articles", fileName)
		f, err := os.Create(articlePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		articleByte, err := json.Marshal(article)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = f.Write(articleByte)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/base.html", "./templates/addArticle.html")
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
	router := http.NewServeMux()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/home", homeHandler)
	router.HandleFunc("/article/{id}", articleHandler)
	router.HandleFunc("/admin/publish", publishHandler)
	router.HandleFunc("/new-article", newArticleHandler)
	http.ListenAndServe(":8080", router)
}
