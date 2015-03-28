package controller

import (
	"net/http"
	"path/filepath"
)

func FileHandler(w http.ResponseWriter, r *http.Request) {
	filename := resourcePath(r.URL.Path)
	if fileExist(filename) {
		http.ServeFile(w, r, filename)
		return
	}

	imgExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if _, ok := imgExts[filepath.Ext(r.URL.Path)]; ok {
		imageResizeHandler(w, r)
		return
	} else {
		http.NotFound(w, r)
	}
}
