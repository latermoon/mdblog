package controller

import (
	"blog"
	"os"
	"path"
	"strings"
	"time"
)

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// return markdown filename
func articlePath(urlpath string) string {
	return path.Join(blog.Path("article"), strings.TrimSuffix(urlpath, path.Ext(urlpath))+".md")
}

func isPrivatePath(urlpath string) bool {
	return strings.Contains(urlpath, "/private/")
}

// Mon, 02 Jan 2006 15:04:05 GMT
func expiresHeader(t time.Time) string {
	loc, _ := time.LoadLocation("GMT")
	return t.In(loc).Format(time.RFC1123)
}
