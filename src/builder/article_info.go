package builder

import (
	"fmt"
	"html/template"
	"time"
)

type ArticleInfo struct {
	Title    string
	Date     time.Time
	Author   string
	Content  template.HTML
	Url      string
	IsPublic bool
	Filename string // full path: ./public/hello.md
	BaseName string // base name: hello.md
	HtmlName string // html name: hello.html
}

func (a *ArticleInfo) DateString() string {
	if a.Date.Unix() < 1 {
		return ""
	}
	return fmt.Sprintf("%d-%d-%d", a.Date.Year(), a.Date.Month(), a.Date.Day())
}

// Sortable
type ArticleInfos []*ArticleInfo

func (a ArticleInfos) Len() int {
	return len(a)
}

func (a ArticleInfos) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ArticleInfos) Less(i, j int) bool {
	return a[i].Date.UnixNano() < a[j].Date.UnixNano()
}
