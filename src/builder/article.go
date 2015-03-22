package builder

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Article struct {
	filename string
	info     os.FileInfo
}

func NewArticle(filename string) (*Article, error) {
	a := &Article{
		filename: filename,
	}
	if info, err := os.Stat(filename); err != nil {
		return nil, err
	} else {
		a.info = info
	}
	return a, nil
}

func (a *Article) Parse() (*ArticleInfo, error) {
	info := &ArticleInfo{}
	info.Filename = a.filename
	info.IsPublic = true // default
	info.Url = a.HtmlName()

	// parse markdown
	md, err := ioutil.ReadFile(a.filename)
	if err != nil {
		return nil, err
	}
	mdhtml := string(blackfriday.MarkdownCommon(md))

	// file info
	a.fillInfo(mdhtml, info)

	return info, nil
}

func (a *Article) fillInfo(mdhtml string, info *ArticleInfo) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(mdhtml))
	if err != nil {
		return err
	}
	// Title
	h1 := doc.Find("h1").First()
	info.Title = h1.Text()
	h1.Remove()

	// Meta
	ul := doc.Find("ul").First()
	ul.Find("li").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.HasPrefix(text, "by") {
			info.Author = text[3:]
		} else if strings.HasPrefix(text, "on") {
			info.Date = text[3:]
		} else if strings.HasPrefix(text, "private") {
			info.IsPublic = false
		}
	})
	ul.Remove()

	if html, err := doc.Html(); err != nil {
		return err
	} else {
		html = strings.TrimPrefix(html, "<html><head></head><body>")
		html = strings.TrimSuffix(html, "</body></html>")
		html = strings.TrimSpace(html)
		info.Content = template.HTML(html)
	}

	return nil
}

func (a *Article) Filename() string {
	return a.filename
}

func (a *Article) BaseName() string {
	_, file := filepath.Split(a.filename)
	return file
}

func (a *Article) HtmlName() string {
	name := a.BaseName()
	comma := strings.LastIndex(name, ",")
	if comma != -1 {
		name = strings.TrimSpace(name[comma+1:])
	}
	ext := filepath.Ext(name)
	return strings.TrimSuffix(name, ext) + ".html"
}

func (a *Article) String() string {
	return a.Filename()
}
