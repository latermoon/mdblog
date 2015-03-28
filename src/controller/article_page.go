package controller

import (
	"blog"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func PublicArticlePage(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(blog.Path("article"), strings.TrimSuffix(r.URL.Path, ".html")+".md")
	art, err := blog.ParseArticle(filename)
	if err != nil {
		log.Println(err)
		ErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("parse error: %s", r.URL.Path))
		return
	}
	renderArticlePage(w, art)
}

func renderArticlePage(wr io.Writer, art *blog.Article) error {
	if err := blog.Template().ExecuteTemplate(wr, "article.tmpl", art); err != nil {
		return err
	}
	return nil
}