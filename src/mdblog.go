package main

import (
	"blog"
	c "controller"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"log"
	"os"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// go run src/mdblog.go /home/workspace/blog.latermoon.me
func main() {
	if len(os.Args) < 2 {
		log.Fatal("workspace not special")
	}

	// init blog
	blog.Init(os.Args[1])

	log.Printf("mdblog sercive start %s\n", blog.Config().Server)
	log.Printf("workspace: %s", blog.Workspace())

	// init martini
	m := blog.Martini()

	// session
	store := sessions.NewCookieStore([]byte(blog.Config().AuthKey))
	store.Options(sessions.Options{Path: "/private/", MaxAge: 24 * 60 * 60})
	m.Use(sessions.Sessions(blog.Config().SessionName, store))

	// static expires
	gmtloc, _ := time.LoadLocation("GMT")
	m.Use(martini.Static(blog.Path("public"), martini.StaticOptions{
		SkipLogging: true,
		Expires:     func() string { return time.Now().In(gmtloc).Add(time.Hour * 24 * 7).Format(time.RFC1123) },
	}))

	// handlers
	m.Use(c.AuthHandler)
	m.Get("/", c.HomePage)
	m.Get("/private/", c.HomePage)
	m.Get(`/(.*).html`, c.ArticlePage)
	m.Post("/login", c.LoginAction)
	m.Get("/logout", c.LogoutAction)
	m.NotFound(c.FileHandler) // include custom resize image

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
