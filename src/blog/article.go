package blog

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type Article struct {
	Title       string        // blog title
	Date        time.Time     // 2015-3-28
	ModTime     time.Time     // modification time
	Author      string        // author
	ContentHTML template.HTML // from markdown
	Filename    string        // article/hello.md
	BaseName    string        // base name: hello.md
	HtmlName    string        // html name: hello.html
	IsPrivate   bool          // private or not
}

func (a *Article) DateString() string {
	dt := a.Date
	if dt.Unix() < 1 {
		return ""
	}
	return fmt.Sprintf("%d-%d-%d", dt.Year(), dt.Month(), dt.Day())
}

// Sortable
type Articles []*Article

func (a Articles) Len() int {
	return len(a)
}

func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Articles) Less(i, j int) bool {
	return a[i].Date.UnixNano() < a[j].Date.UnixNano()
}

// parse markdown article in to Article object
func parseArticle(filename string) (*Article, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	art := &Article{}
	art.Filename = filename
	art.BaseName = path.Base(filename)
	art.HtmlName = htmlName(art.BaseName)
	art.Date = time.Unix(0, 0)
	art.ModTime = info.ModTime()

	// parse markdown
	md, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	mdhtml := string(blackfriday.MarkdownCommon(md))

	// file info
	if err := fillArticleInfo(mdhtml, art); err != nil {
		return nil, err
	}
	return art, nil
}

func fillArticleInfo(mdhtml string, art *Article) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(mdhtml))
	if err != nil {
		return err
	}
	// Title
	h1 := doc.Find("h1").First()
	art.Title = h1.Text()
	h1.Remove()

	// Meta
	ul := doc.Find("ul").First()
	found := false
	ul.Find("li").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.HasPrefix(text, "by") {
			art.Author, found = text[3:], true
		} else if strings.HasPrefix(text, "on") {
			art.Date, found = parseDate(text[3:]), true
		}
	})
	if found {
		ul.Remove()
	}

	if html, err := doc.Html(); err != nil {
		return err
	} else {
		html = strings.TrimPrefix(html, "<html><head></head><body>")
		html = strings.TrimSuffix(html, "</body></html>")
		html = strings.TrimSpace(html)
		art.ContentHTML = template.HTML(html)
	}

	return nil
}
