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
	readerHomePage(w, arts)
}

func readerHomePage(wr io.Writer, lst []*blog.Article) error {
	data := map[string]interface{}{
		"Title":    blog.Config().Title,
		"Articles": lst,
	}
	if err := blog.Template().ExecuteTemplate(wr, "home.tmpl", data); err != nil {
		return err
	}
	return nil
}

func PrivateHomePage(w http.ResponseWriter, r *http.Request) {

}
