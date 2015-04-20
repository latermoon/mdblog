package controller

import (
	"github.com/latermoon/mdblog/blog"
	"log"
	"net/http"
	"path"
	"sort"
	"strings"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	dirname, title := path.Join(blog.Path("article"), r.URL.Path), blog.Config().Title
	isPrivate := strings.Contains(r.URL.Path, "/private/")
	if isPrivate {
		title = "Private Blog"
	}

	arts, err := blog.GetAllArticles(dirname)
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
