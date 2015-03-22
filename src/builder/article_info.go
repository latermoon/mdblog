package builder

import (
	"html/template"
)

type ArticleInfo struct {
	Title    string
	Date     string
	Author   string
	Content  template.HTML
	Url      string
	IsPublic bool
	Filename string
}
