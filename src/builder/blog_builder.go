package builder

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type BlogBuilder struct {
	articles    []*Article
	workspace   string
	articlePath string
	publicPath  string
	template    *template.Template
}

func NewBlogBuilder(workspace string) (*BlogBuilder, error) {
	builder := &BlogBuilder{
		workspace: workspace,
	}
	builder.articlePath = filepath.Join(workspace, "article")
	builder.publicPath = filepath.Join(workspace, "public")
	return builder, nil
}

func (b *BlogBuilder) RebuildAll() {
	if err := b.initTemplate(); err != nil {
		log.Fatal(err)
	}

	if err := b.reloadArticles(); err != nil {
		log.Fatal(err)
	}

	infos := make([]*ArticleInfo, 0)
	for _, article := range b.articles {
		info, err := article.Parse()
		if err != nil {
			log.Fatal(err)
		}

		htmlfile := filepath.Join(b.publicPath, article.HtmlName())

		if info.IsPublic {
			infos = append(infos, info)
			log.Println("parse", article.BaseName())
			if err := b.renderArticle(info, htmlfile); err != nil {
				log.Fatal(err)
			}
		} else {
			if _, err := os.Stat(htmlfile); err == nil {
				log.Printf("remove prvate blog %s\n", article.HtmlName())
				os.Remove(htmlfile)
			} else {
				log.Printf("skip private blog %s\n", article.BaseName())
			}
		}
	}

	homeInfo := &HomeInfo{
		Title:    "latermoon's blog",
		Articles: infos,
	}
	indexName := filepath.Join(b.publicPath, "index.html")
	log.Println("rebuild index.html with", len(infos), "articles")
	if err := b.renderHome(homeInfo, indexName); err != nil {
		log.Fatal(err)
	}
}

func (b *BlogBuilder) renderHome(info *HomeInfo, filename string) error {
	buf := &bytes.Buffer{}
	if err := b.template.ExecuteTemplate(buf, "home.tmpl", info); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, buf.Bytes(), 0666)
}

func (b *BlogBuilder) renderArticle(info *ArticleInfo, filename string) error {
	buf := &bytes.Buffer{}
	if err := b.template.ExecuteTemplate(buf, "article.tmpl", info); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, buf.Bytes(), 0666)
}

func (b *BlogBuilder) reloadArticles() error {
	b.articles = make([]*Article, 0)
	err := filepath.Walk(b.articlePath, func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		article, err := NewArticle(filename)
		if err != nil {
			return err
		}
		b.articles = append(b.articles, article)
		return nil
	})
	sort.Sort(sort.Reverse(Articles(b.articles)))
	return err
}

func (b *BlogBuilder) initTemplate() error {
	tmpl, err := template.ParseFiles(
		filepath.Join(b.workspace, "template", "article.tmpl"),
		filepath.Join(b.workspace, "template", "home.tmpl"),
	)
	b.template = tmpl
	return err
}
