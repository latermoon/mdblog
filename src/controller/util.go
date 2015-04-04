package controller

import (
	"blog"
	"crypto/md5"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func encodeMd5Password(md5pwd string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(blog.Config().Salt+md5pwd)))
}

func encodeRawPassword(pwd string) string {
	md5pwd := fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
	return encodeMd5Password(md5pwd)
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
