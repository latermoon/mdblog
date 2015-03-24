package builder

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Article struct {
	filename string
	info     *ArticleInfo
}

func NewArticle(filename string) (*Article, error) {
	a := &Article{
		filename: filename,
	}
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Article) Parse() (*ArticleInfo, error) {
	info := &ArticleInfo{}
	info.Filename = a.filename
	info.IsPublic = true // default
	info.Url = a.HtmlName()
	info.Date = time.Unix(0, 0)

	// parse markdown
	md, err := ioutil.ReadFile(a.filename)
	if err != nil {
		return nil, err
	}
	mdhtml := string(blackfriday.MarkdownCommon(md))

	// file info
	a.fillInfo(mdhtml, info)

	a.info = info
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

func parseDate(s string) time.Time {
	re := regexp.MustCompile(`(\d{4})\-(\d{1,2})-(\d{1,2})`)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 0 {
		return time.Unix(0, 0)
	}
	year, e1 := strconv.Atoi(matches[1])
	month, e2 := strconv.Atoi(matches[2])
	day, e3 := strconv.Atoi(matches[3])
	if e1 != nil || e2 != nil || e3 != nil {
		return time.Unix(0, 0)
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

func (a *Article) Info() *ArticleInfo {
	return a.info
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
