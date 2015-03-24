package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	"net/http"
	"path/filepath"
	"strings"
)

func privateGroup(r martini.Router) {
	r.Get(`/(.*).html`, privateArticleHandler)
	r.Get(`/(.*)`, privateFileHandler)
}

func authHandler() martini.Handler {
	return auth.BasicFunc(func(username, password string) bool {
		return username == "latermoon" && password == "1234"
		// t := md5.New()
		// t.Write([]byte(password))
		// encpwd := fmt.Sprintf("%x", t.Sum(nil))
		// return username == "yan" && encpwd == "8472bd8ee3d641fd1225cbd289075a33"
	})
}

func privateFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(Workspace, r.URL.Path)
	http.ServeFile(w, r, filename)
}

func privateArticleHandler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(Workspace, strings.TrimSuffix(r.URL.Path, ".html")+".md")
	serveArticle(w, r, filename)
}
