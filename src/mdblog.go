package main

import (
	"blog"
	c "controller"
	"log"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// go run src/mdblog.go /home/workspace/blog.latermoon.me
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
	m.NotFound(c.Static("article"), c.FileHandler)

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
