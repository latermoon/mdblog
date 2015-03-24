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
	m.Use(sessions.Sessions("sess", sessions.NewCookieStore([]byte("auth"))))
	m.Use(martini.Static(filepath.Join(workspace, "public")))
	m.Get("/", publicIndexHandler)
	m.Get(`/([^\/]*).html`, publicArticleHandler)
	m.Get("/img/(.*)", imageResizeHandler)
	m.Post("/auth", authHandler)
	m.Group("/private", privateGroup)
	m.Get("/private/img/(.*)", imageResizeHandler)
	m.RunOnAddr(addr)

	<-done
}
