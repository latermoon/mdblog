package server

import (
	"builder"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"log"
	"path/filepath"
)

var Workspace string // website home directory
var blogBuilder *builder.BlogBuilder

const sessName = "auth2"

func ListenAndServe(addr string, workspace string) {
	Workspace = workspace

	// builder
	var err error
	blogBuilder, err = builder.NewBlogBuilder(Workspace)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool)

	log.Println("workspace:", workspace)
	// go watch(filepath.Join(workspace, "article"), filepath.Join(workspace, "template"))

	m := martini.Classic()
	store := sessions.NewCookieStore([]byte(sessName))
	store.Options(sessions.Options{
		Path:   "/private/",
		MaxAge: 24 * 60 * 60, // one day
	})
	m.Use(sessions.Sessions("sess", store))
	m.Use(martini.Static(filepath.Join(workspace, "public")))
	m.Get("/", publicIndexHandler)
	m.Get(`/([^\/]*).html`, publicArticleHandler)
	m.Get("/img/(.*)", imageResizeHandler)
	m.Post("/auth", authHandler)
	m.Get("/logout", logoutHandler)
	m.Group("/private", privateGroup)
	m.RunOnAddr(addr)

	<-done
}
