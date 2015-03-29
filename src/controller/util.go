package controller

import (
	"blog"
	"os"
	"path"
	"strings"
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
