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
	builder.initTemplate()
	return builder, nil
}

// Rebuild All
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
			log.Printf("skip private %s\n", article.BaseName())
		}
	}

	b.removeUnuseHtml()

	sort.Sort(sort.Reverse(ArticleInfos(infos)))

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

// when you set a article as `- private` or delete it, it's html should be remove from public
func (b *BlogBuilder) removeUnuseHtml() {
	// prevent errors caused by an empty array bug
	if len(b.articles) == 0 {
		return
	}
	filesInuse := make(map[string]bool)
	for _, art := range b.articles {
		filesInuse[art.HtmlName()] = true
	}

	files, err := ioutil.ReadDir(b.publicPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		} else if file.Name() == "index.html" {
			continue
		}
		if _, ok := filesInuse[file.Name()]; !ok {
			log.Println("remove unuse ", file.Name())
			os.Remove(filepath.Join(b.publicPath, file.Name()))
		}
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
	files, err := ioutil.ReadDir(b.articlePath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		art, err := NewArticle(filepath.Join(b.articlePath, file.Name()))
		if err != nil {
			return err
		}
		b.articles = append(b.articles, art)
	}
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

func (b *BlogBuilder) Template() *template.Template {
	return b.template
}
