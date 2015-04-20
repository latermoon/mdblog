package blog

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

var _cachelock sync.Mutex

func GetArticle(filename string) (*Article, error) {
	_cachelock.Lock()
	defer _cachelock.Unlock()

	key := "article:" + filename
	if val, ok := _caches.Get(key); ok {
		if info, err := os.Stat(filename); err != nil {
			return nil, err
		} else if info.ModTime() != val.(*Article).ModTime {
			return reloadArticleCache(filename)
		} else {
			return val.(*Article), nil
		}
	} else {
		return reloadArticleCache(filename)
	}
}

func reloadArticleCache(filename string) (*Article, error) {
	log.Println("ParseArticle:", filename)
	art, err := parseArticle(filename)
	if err != nil {
		return nil, err
	}
	key := "article:" + filename
	_caches.Add(key, art)
	return art, nil
}

func GetAllArticles(dirname string) ([]*Article, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	lst := make([]*Article, 0, 100)
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		art, err := GetArticle(path.Join(dirname, file.Name()))
		if err != nil {
			return nil, err
		}
		lst = append(lst, art)
	}
	return lst, nil
}
