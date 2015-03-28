package server

import (
	"blog"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func privateGroup(r martini.Router) {
	r.Get("/", privateIndexHandler)
	r.Get(`/(.*).html`, privateArticleHandler)
	r.Get("/password.txt", func(w http.ResponseWriter) { w.WriteHeader(http.StatusForbidden) })
	r.Get(`/(.*)`, privateFileHandler)
}

func privateIndexHandler(w http.ResponseWriter, r *http.Request, session sessions.Session) {
	dirname := blog.Path("private")
	if ok := checkAuth(w, r, session); !ok {
		return
	}
	serveIndex(w, r, false, dirname)
}

// <script>var pwd = prompt('Your password?');alert(pwd);</script>
func privateArticleHandler(w http.ResponseWriter, r *http.Request, session sessions.Session) {
	filename := filepath.Join(blog.Workspace(), strings.TrimSuffix(r.URL.Path, ".html")+".md")
	if ok := checkAuth(w, r, session); !ok {
		return
	}
	serveArticle(w, r, false, filename)
}

func privateFileHandler(w http.ResponseWriter, r *http.Request, session sessions.Session) {
	filename := filepath.Join(blog.Workspace(), r.URL.Path)
	if ok := checkAuth(w, r, session); !ok {
		return
	}
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
	auth := session.Get(sessName)
	if auth != salt {
		io.WriteString(w, authFormStirng)
		return false
	}
	return true
}
