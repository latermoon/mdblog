package controller

import (
	"blog"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

func PrivateGroup(r martini.Router) {
	r.Get("/", PrivateHomePage)
	r.Get(`/(.*).html`, PrivateArticlePage)
	r.Get("/password.txt", func(w http.ResponseWriter) { w.WriteHeader(http.StatusForbidden) })
	r.Get(`/(.*)`, privateFileHandler)
}

func PrivateHomePage(w http.ResponseWriter, r *http.Request, session sessions.Session) {
	if ok := checkAuth(w, r, session); !ok {
		return
	}
	arts, err := blog.ParseAllArticles(blog.Path("private"))
	if err != nil {
		log.Println(err)
		ErrorPage(w, http.StatusInternalServerError, "parse error")
		return
	}

	sort.Sort(sort.Reverse(blog.Articles(arts)))

	data := map[string]interface{}{
		"Title":    "Private Blog",
		"Articles": arts,
	}
	readerHomePage(w, data)
}

func privateFileHandler(w http.ResponseWriter, r *http.Request, session sessions.Session) {
	if ok := checkAuth(w, r, session); !ok {
		return
	}
	filename := filepath.Join(blog.Workspace(), r.URL.Path)
	if _, err := os.Stat(filename); err == nil {
		http.ServeFile(w, r, filename)
		return
	}

	imgExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if _, ok := imgExts[filepath.Ext(r.URL.Path)]; ok {
		imageResizeHandler(w, r)
		return
	}
}

func checkAuth(w http.ResponseWriter, r *http.Request, session sessions.Session) bool {
	salt := blog.Config().Password
	auth := session.Get(blog.Config().AuthKey)
	if auth != salt {
		io.WriteString(w, authFormStirng)
		return false
	}
	return true
}
