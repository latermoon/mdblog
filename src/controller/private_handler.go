package controller

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"net/http"
)

func PrivateGroup(r martini.Router) {
	r.Get("/", PrivateHomePage)
	r.Get(`/(.*).html`, PrivateArticlePage)
	r.Get("/password.txt", func(w http.ResponseWriter) { w.WriteHeader(http.StatusForbidden) })
	r.Get(`/(.*)`, privateFileHandler)
}

func PrivateArticlePage(w http.ResponseWriter, r *http.Request, session sessions.Session) {

}

func privateFileHandler(w http.ResponseWriter, r *http.Request, session sessions.Session) {

}
