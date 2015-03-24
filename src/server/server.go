package server

import (
	"builder"
	"github.com/go-martini/martini"
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

	log.Println("watching:", workspace)
	go watch(filepath.Join(workspace, "article"), filepath.Join(workspace, "template"))

	m := martini.Classic()
	m.Use(martini.Static(filepath.Join(workspace, "public")))
	m.Get(`/([^\/]*).html`, publicArticleHandler)
	m.Get("/img/(.*)", imageResizeHandler)
	m.Group("/private", privateGroup)
	m.RunOnAddr(addr)

	<-done
}
