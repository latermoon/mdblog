package server

import (
	"blog"
	"github.com/go-martini/martini"
	"github.com/hashicorp/golang-lru"
	"github.com/howeyc/fsnotify"
	"github.com/martini-contrib/sessions"
	"html/template"
	"log"
	"path/filepath"
	"time"
)

var templates *template.Template
var cache *lru.Cache

const (
	PublicIndexCache  = "_index"
	PrivateIndexCache = "_private"
)

const sessName = "auth"

func ListenAndServe(addr string, dir string) {
	log.Println("workspace:", blog.Workspace())
	if err := blog.Init(dir); err != nil {
		log.Fatal(err)
	}

	// init template
	if err := initTemplate(); err != nil {
		log.Fatal(err)
	}

	cache, _ = lru.New(100)

	// watch template directory
	watchTemplateModify()
	watchAtricleModify()

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
	go watch(blog.Path("template"), func(e *fsnotify.FileEvent, err error) {
		log.Println("event:", e.String())
		cache.Purge()
		initTemplate()
	})
}

func watchAtricleModify() {
	go watch(blog.Path("article"), func(e *fsnotify.FileEvent, err error) {
		log.Println("event:", e.String())
		cache.Remove(PublicIndexCache)
	})

	go watch(blog.Path("private"), func(e *fsnotify.FileEvent, err error) {
		log.Println("event:", e.String())
		cache.Remove(PrivateIndexCache)
	})
}
