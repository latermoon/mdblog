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
	Filename string
}

func (a *ArticleInfo) DateString() string {
	if a.Date.Unix() < 1 {
		return ""
	}
	return fmt.Sprintf("%d-%d-%d", a.Date.Year(), a.Date.Month(), a.Date.Day())
}
