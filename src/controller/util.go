package controller

import (
	"blog"
	"os"
	"path/filepath"
	"strings"
)

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// return filename and private or not
func articlePath(urlpath string) (string, bool) {
	if !strings.HasPrefix(urlpath, "/private/") {
		return filepath.Join(blog.Path("article"), strings.TrimSuffix(urlpath, ".html")+".md"), false
	} else {
		return filepath.Join(blog.Workspace(), strings.TrimSuffix(urlpath, ".html")+".md"), true
	}
}

func resourcePath(urlpath string) string {
	if !strings.HasPrefix(urlpath, "/private/") {
		return filepath.Join(blog.Path("public"), urlpath)
	} else {
		return filepath.Join(blog.Workspace(), urlpath)
	}
}
