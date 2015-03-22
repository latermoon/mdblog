package builder

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

	tmpl, err := template.ParseFiles(
		filepath.Join(workspace, "template", "article.tmpl"),
		filepath.Join(workspace, "template", "home.tmpl"),
	)
	if err != nil {
		return nil, err
	}
	builder.template = tmpl
	builder.articlePath = filepath.Join(workspace, "article")
	builder.publicPath = filepath.Join(workspace, "public")
	return builder, nil
}

func (b *BlogBuilder) RebuildAll() error {
	b.reloadArticles()

	infos := make([]*ArticleInfo, 0)
	for _, article := range b.articles {
		log.Println("parse", article.Filename())
		info, err := article.Parse()
		if err != nil {
			return err
		}
		infos = append(infos, info)

		filename := filepath.Join(b.publicPath, info.Filename)
		if err := b.renderArticle(info, filename); err != nil {
			return err
		}
	}

	homeInfo := &HomeInfo{
		Title:    "latermoon's blog",
		Articles: infos,
	}
	indexName := filepath.Join(b.publicPath, "index.html")
	if err := b.renderHome(homeInfo, indexName); err != nil {
		return err
	}
	return nil
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

func (b *BlogBuilder) reloadArticles() {
	b.articles = make([]*Article, 0)
	filepath.Walk(b.articlePath, func(filename string, info os.FileInfo, err error) error {
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
}
