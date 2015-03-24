package server

import (
	"builder"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

func publicArticleHandler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(Workspace, "article", strings.TrimSuffix(r.URL.Path, ".html")+".md")
	serveArticle(w, r, filename)
}

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
