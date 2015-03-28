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

func articlePath(urlpath string) string {
	if !strings.HasPrefix(urlpath, "/private/") {
		return filepath.Join(blog.Path("article"), strings.TrimSuffix(urlpath, ".html")+".md")
	} else {
		return filepath.Join(blog.Workspace(), strings.TrimSuffix(urlpath, ".html")+".md")
	}
}

func resourcePath(urlpath string) string {
	if !strings.HasPrefix(urlpath, "/private/") {
		return filepath.Join(blog.Path("public"), urlpath)
	} else {
		return filepath.Join(blog.Workspace(), urlpath)
	}
}
