package server

import (
	"github.com/go-martini/martini"
	"github.com/howeyc/fsnotify"
	"github.com/martini-contrib/sessions"
	"html/template"
	"log"
	"path/filepath"
	"time"
)

var Workspace string // website home directory
var templates *template.Template
var blogConfig *BlogConfig

const sessName = "auth"

func ListenAndServe(addr string, dir string) {
	Workspace = dir
	log.Println("workspace:", dir)

	// init template
	if err := initTemplate(); err != nil {
		log.Fatal(err)
	}

	// init blog config
	if cfg, err := NewBlogConfig(filepath.Join(Workspace, "blog.txt")); err != nil {
		log.Fatal(err)
	} else {
		blogConfig = cfg
	}
	log.Println(blogConfig)

	// watch template directory
	go watchTemplateModify()

	// http server
	m := martini.Classic()
	store := sessions.NewCookieStore([]byte(sessName))
	store.Options(sessions.Options{
		Path:   "/private/",
		MaxAge: 24 * 60 * 60, // one day
	})
	m.Use(sessions.Sessions("sess", store))
	gmtLoc, _ := time.LoadLocation("GMT")
	m.Use(martini.Static(filepath.Join(dir, "public"), martini.StaticOptions{
		SkipLogging: true,
		Expires:     func() string { return time.Now().In(gmtLoc).Add(time.Hour * 24 * 7).Format(time.RFC1123) },
	}))
	m.Get("/", publicIndexHandler)
	m.Group("/private", privateGroup)
	m.Get(`/(.*).html`, publicArticleHandler)
	m.Get("/img/(.*)", imageResizeHandler)
	m.Post("/auth", authHandler)
	m.Get("/logout", logoutHandler)

	m.RunOnAddr(addr)
}

func watchTemplateModify() {
	dir := filepath.Join(Workspace, "template")
	watch(dir, func(e *fsnotify.FileEvent, err error) {
		log.Println("event:", e.String())
		initTemplate()
	})
}
