package builder

import (
	"os"
)

type Article struct {
	filename string
	info     os.FileInfo
}

func NewArticle(filename string) (*Article, error) {
	a := &Article{
		filename: filename,
	}
	if info, err := os.Stat(filename); err != nil {
		return nil, err
	} else {
		a.info = info
	}
	return a, nil
}

func (a *Article) Parse() (*ArticleInfo, error) {
	return nil, nil
}

func (a *Article) Filename() string {
	return a.info.Name()
}

func (a *Article) String() string {
	return a.Filename()
}
