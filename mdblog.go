package main

import (
	"github.com/go-martini/martini"
	"github.com/latermoon/mdblog/blog"
	c "github.com/latermoon/mdblog/controller"
	"log"
	"net/http"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	martini.Env = martini.Prod
}

// go run mdblog.go /home/workspace/blog.latermoon.me
func main() {
	if len(os.Args) < 2 {
		log.Fatal("workspace not special")
	}

	// init Workspace first
	blog.Init(os.Args[1])

	log.Printf("mdblog sercive start %s\n", blog.Config().Server)
	log.Printf("workspace: %s", blog.Workspace())

	// martini
	m := blog.Martini()

	// middlewares
	m.Use(c.Sessions())
	m.Use(c.Static("static"))

	// handlers
	m.Use(c.AuthHandler)
	m.Get("(.*)/", c.HomePage)
	m.Get("/(.*).(html|md)", c.ArticlePage)
	m.Post("/login", c.LoginAction)
	m.Get("/logout", c.LogoutAction)
	m.NotFound(c.Static("article"), c.ImageResize("article"), http.NotFound)

	// Go!
	m.RunOnAddr(blog.Config().Server)
}

/*
	go watch(blog.Path("template"), func(e *fsnotify.FileEvent, err error) {
		log.Println("event:", e.String())
		cache.Purge()
		initTemplate()
	})
*/
