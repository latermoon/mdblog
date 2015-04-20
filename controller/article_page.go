package controller

import (
	"github.com/latermoon/mdblog/blog"
	"log"
	"net/http"
)

func ArticlePage(w http.ResponseWriter, r *http.Request) {
	filename := articlePath(r.URL.Path)
	if !fileExist(filename) {
		http.NotFound(w, r)
		return
	}

	art, err := blog.GetArticle(filename)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	art.IsPrivate = isPrivatePath(r.URL.Path)

	if err := blog.Template().ExecuteTemplate(w, "article.tmpl", art); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
