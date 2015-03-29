package controller

import (
	"blog"
	"github.com/go-martini/martini"
	"time"
)

func Static(dirname string) martini.Handler {
	return martini.Static(blog.Path(dirname), martini.StaticOptions{
		SkipLogging: true,
		Expires:     func() string { return expiresHeader(time.Now().Add(time.Hour * 24 * 7)) },
	})
}
