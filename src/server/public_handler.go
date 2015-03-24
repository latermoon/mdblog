package server

import (
	"builder"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
)

func publicArticleHandler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(Workspace, "article", strings.TrimSuffix(r.URL.Path, ".html")+".md")
	serveArticle(w, r, filename)
}

func publicIndexHandler(w http.ResponseWriter, r *http.Request) {
	dirname := filepath.Join(Workspace, "article")
	serveIndex(w, r, dirname)
}

// render markdown article to html
func serveArticle(w http.ResponseWriter, r *http.Request, filename string) {
	parser := builder.NewArticleParser(filename)
	info, err := parser.Parse()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "parse error: %s", err)
		return
	}

	if err := blogBuilder.Template().ExecuteTemplate(w, "article.tmpl", info); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "parse error: %s", err)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request, dirname string) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "parse error: %s", err)
		return
	}
	infos := make([]*builder.ArticleInfo, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		parser := builder.NewArticleParser(filepath.Join(dirname, file.Name()))
		info, err := parser.Parse()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "parse error: %s", err)
			return
		}

		infos = append(infos, info)
	}
	sort.Sort(sort.Reverse(builder.ArticleInfos(infos)))
	data := map[string]interface{}{"Articles": infos, "Title": "latermoon's blog"}
	if err := blogBuilder.Template().ExecuteTemplate(w, "home.tmpl", data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "parse error: %s", err)
		return
	}
}
