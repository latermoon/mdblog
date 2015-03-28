package controller

import (
	"blog"
	"log"
	"net/http"
	"sort"
	"strings"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	var dirname, title string
	isPrivate := false
	if strings.HasPrefix(r.URL.Path, "/private/") {
		dirname, title = blog.Path("private"), "Private Blog"
		isPrivate = true
	} else {
		dirname, title = blog.Path("article"), blog.Config().Title
	}

	arts, err := blog.ParseAllArticles(dirname)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sort.Sort(sort.Reverse(blog.Articles(arts)))

	data := map[string]interface{}{
		"Title":     title,
		"IsPrivate": isPrivate,
		"Articles":  arts,
	}

	if err := blog.Template().ExecuteTemplate(w, "home.tmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
