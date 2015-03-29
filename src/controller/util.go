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

// return markdown filename
func articlePath(urlpath string) string {
	return filepath.Join(blog.Path("article"), strings.TrimSuffix(urlpath, filepath.Ext(urlpath))+".md")
}

func isPrivatePath(urlpath string) bool {
	return strings.Contains(urlpath, "/private/")
}

func resourcePath(urlpath string) string {
	if !strings.HasPrefix(urlpath, "/private/") {
		return filepath.Join(blog.Path("public"), urlpath)
	} else {
		return filepath.Join(blog.Workspace(), urlpath)
	}
}
