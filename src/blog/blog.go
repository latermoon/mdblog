package blog

import (
	"github.com/BurntSushi/toml"
	"github.com/go-martini/martini"
	"html/template"
	"path/filepath"
)

var _workspace string
var _config BlogConfig
var _templates *template.Template
var _martini = martini.Classic()

// Must call blog.Init(...) first
func Init(workspace string) error {
	_workspace = workspace
	var err error

	// load config
	if _, err = toml.DecodeFile(Path("blog.txt"), &_config); err != nil {
		return err
	}
	_config.SessionName = "sess"
	_config.AuthKey = "auth"

	// init templates
	if _templates, err = template.ParseFiles(Path("template/article.tmpl"), Path("template/home.tmpl")); err != nil {
		return err
	}

	return nil
}

func Martini() *martini.ClassicMartini {
	return _martini
}

func Config() BlogConfig {
	return _config
}

func Template() *template.Template {
	return _templates
}

func Workspace() string {
	return _workspace
}

func Path(dir string) string {
	return filepath.Join(Workspace(), dir)
}
