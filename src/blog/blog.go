package blog

import (
	"github.com/BurntSushi/toml"
	"html/template"
	"path/filepath"
)

var _workspace string
var _config BlogConfig
var _templates *template.Template

// Must call blog.Init(...) first
func Init(workspace string) error {
	_workspace = workspace
	var err error

	// load config
	if _, err = toml.DecodeFile(Path("blog.txt"), &_config); err != nil {
		return err
	}

	// init templates
	if _templates, err = template.ParseFiles(Path("template/article.tmpl"), Path("template/home.tmpl")); err != nil {
		return err
	}

	return nil
}

func Workspace() string {
	return _workspace
}

func Path(dir string) string {
	return filepath.Join(Workspace(), dir)
}

func Config() BlogConfig {
	return _config
}

func Template() *template.Template {
	return _templates
}
