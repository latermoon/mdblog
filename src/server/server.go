package server

import (
	"github.com/go-martini/martini"
	"github.com/howeyc/fsnotify"
	"github.com/martini-contrib/sessions"
	"html/template"
	"log"
	"path/filepath"
)

var Workspace string // website home directory
var templates *template.Template

const sessName = "auth"

func ListenAndServe(addr string, dir string) {
	Workspace = dir

	if err := initTemplate(); err != nil {
		log.Fatal(err)
	}

	log.Println("workspace:", dir)
	go watchTemplateModify()

	m := martini.Classic()
	store := sessions.NewCookieStore([]byte(sessName))
	store.Options(sessions.Options{
		Path:   "/private/",
		MaxAge: 24 * 60 * 60, // one day
	})
	m.Use(sessions.Sessions("sess", store))
	m.Use(martini.Static(filepath.Join(dir, "public")))
	m.Get("/", publicIndexHandler)
	m.Get(`/([^\/]*).html`, publicArticleHandler)
	m.Get("/img/(.*)", imageResizeHandler)
	m.Post("/auth", authHandler)
	m.Get("/logout", logoutHandler)
	m.Group("/private", privateGroup)
	m.RunOnAddr(addr)
}

func watchTemplateModify() {
	dir := filepath.Join(Workspace, "template")
	watch(dir, func(e *fsnotify.FileEvent, err error) {
		log.Println("event:", e.String())
		initTemplate()
	})
}
