package builder

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type ArticleParser struct {
	filename   string
	htmlResult string
}

func NewArticleParser(filename string) *ArticleParser {
	return &ArticleParser{
		filename: filename,
	}
}

func (a *ArticleParser) Parse() (*ArticleInfo, error) {
	if _, err := os.Stat(a.filename); err != nil {
		return nil, err
	}
	info := &ArticleInfo{}
	info.IsPublic = true // default
	info.Filename = a.filename
	info.BaseName = baseName(a.filename)
	info.HtmlName = htmlName(info.BaseName)
	info.Url = info.HtmlName
	info.Date = time.Unix(0, 0)

	// parse markdown
	md, err := ioutil.ReadFile(a.filename)
	if err != nil {
		return nil, err
	}
	mdhtml := string(blackfriday.MarkdownCommon(md))

	// file info
	if err := a.fillInfo(mdhtml, info); err != nil {
		return nil, err
	}
	return info, nil
}

func (a *ArticleParser) fillInfo(mdhtml string, info *ArticleInfo) error {
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
			info.Date = parseDate(text[3:])
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
