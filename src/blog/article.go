package blog

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Article struct {
	Title        string        // blog title
	Date         time.Time     // 2015-3-28
	LastModified time.Time     // file modification time
	Author       string        // author
	ContentHTML  template.HTML // from markdown
	Filename     string        // article/hello.md
	BaseName     string        // base name: hello.md
	HtmlName     string        // html name: hello.html
}

func (a *Article) DateString() string {
	dt := a.Date
	if dt.Unix() < 1 {
		return ""
	}
	return fmt.Sprintf("%d-%d-%d", dt.Year(), dt.Month(), dt.Day())
}

func ParseAllArticles(dirname string) ([]*Article, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	lst := make([]*Article, 0, 100)
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		art, err := ParseArticle(filepath.Join(dirname, file.Name()))
		if err != nil {
			return nil, err
		}
		lst = append(lst, art)
	}
	return lst, nil
}

// parse markdown article in to Article object
func ParseArticle(filename string) (*Article, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}
	art := &Article{}
	art.Filename = filename
	art.BaseName = baseName(filename)
	art.HtmlName = htmlName(art.BaseName)
	art.Date = time.Unix(0, 0)

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
	ul.Find("li").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.HasPrefix(text, "by") {
			art.Author = text[3:]
		} else if strings.HasPrefix(text, "on") {
			art.Date = parseDate(text[3:])
		}
	})
	ul.Remove()

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
