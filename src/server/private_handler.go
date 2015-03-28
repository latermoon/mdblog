package server

import (
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
	dirname := filepath.Join(Workspace, "private")
	if ok := checkAuth(w, r, session); !ok {
		return
	}
	serveIndex(w, r, false, dirname)
}

// <script>var pwd = prompt('Your password?');alert(pwd);</script>
func privateArticleHandler(w http.ResponseWriter, r *http.Request, session sessions.Session) {
	filename := filepath.Join(Workspace, strings.TrimSuffix(r.URL.Path, ".html")+".md")
	if ok := checkAuth(w, r, session); !ok {
		return
	}
	serveArticle(w, r, false, filename)
}

func privateFileHandler(w http.ResponseWriter, r *http.Request, session sessions.Session) {
	filename := filepath.Join(Workspace, r.URL.Path)
	if ok := checkAuth(w, r, session); !ok {
		return
	}
	if _, err := os.Stat(filename); err != nil && strings.HasPrefix(r.URL.Path, "/private/img/") {
		imageResizeHandler(w, r)
		return
	}
	http.ServeFile(w, r, filename)
}

func checkAuth(w http.ResponseWriter, r *http.Request, session sessions.Session) bool {
	salt := blogConfig.Password
	auth := session.Get(sessName)
	if auth != salt {
		io.WriteString(w, authFormStirng)
		return false
	}
	return true
}
