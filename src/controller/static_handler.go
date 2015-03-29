package controller

import (
	"blog"
	"github.com/go-martini/martini"
	"time"
)

func Static(dirname string) martini.Handler {
	gmtloc, _ := time.LoadLocation("GMT")
	return martini.Static(blog.Path(dirname), martini.StaticOptions{
		SkipLogging: true,
		Expires:     func() string { return time.Now().In(gmtloc).Add(time.Hour * 24 * 7).Format(time.RFC1123) },
	})
}
