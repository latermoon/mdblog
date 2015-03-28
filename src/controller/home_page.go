package controller

import (
	"blog"
	"io"
	"log"
	"net/http"
	"sort"
)

func PublicHomePage(w http.ResponseWriter, r *http.Request) {
	arts, err := blog.ParseAllArticles(blog.Path("article"))
	if err != nil {
		log.Println(err)
		ErrorPage(w, http.StatusInternalServerError, "parse error")
		return
	}

	sort.Sort(sort.Reverse(blog.Articles(arts)))

	data := map[string]interface{}{
		"Title":    blog.Config().Title,
		"Articles": arts,
	}
	readerHomePage(w, data)
}

func readerHomePage(wr io.Writer, data map[string]interface{}) error {
	if err := blog.Template().ExecuteTemplate(wr, "home.tmpl", data); err != nil {
		return err
	}
	return nil
}
