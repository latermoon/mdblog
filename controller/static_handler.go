package controller

import (
	"github.com/go-martini/martini"
	"github.com/latermoon/mdblog/blog"
	"time"
)

func Static(dirname string) martini.Handler {
	return martini.Static(blog.Path(dirname), martini.StaticOptions{
		SkipLogging: true,
		Expires:     func() string { return expiresHeader(time.Now().Add(time.Hour * 24 * 7)) },
	})
}
