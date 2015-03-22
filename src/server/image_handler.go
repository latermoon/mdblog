package server

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func imageHandler(w http.ResponseWriter, r *http.Request) {
	dir, _, srcname, sizes := fileInfo(r.URL.Path)
	srcfile := filepath.Join(Workspace, "public", dir, srcname)
	log.Println(r.URL.Path, srcname, sizes)
	// expired
	w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d", 3600))
	// w.Header().Add("Content-Length", "46639")
	// w.Header().Add("Last-Modified", "Sat, 21 Mar 2015 18:33:18 GMT")
	// w.Header().Add("Expires", "Sat, 28 Mar 2016 13:40:31 GMT")
	mimg, err := loadImage(srcfile)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	rimg := resize.Thumbnail(uint(sizes[0]), uint(sizes[1]), mimg, resize.Lanczos3)
	jpeg.Encode(w, rimg, &jpeg.Options{Quality: 75})
}

func loadImage(filename string) (img image.Image, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ext := filepath.Ext(filename)
	switch strings.ToLower(ext) {
	case ".jpg":
		return jpeg.Decode(f)
	case ".jpeg":
		return jpeg.Decode(f)
	case ".png":
		return png.Decode(f)
	default:
		img, _, err = image.Decode(f)
		return
	}
}

func fileInfo(filename string) (dir, name, srcname string, sizes []int) {
	dir, name = filepath.Split(filename)
	ext := filepath.Ext(name)
	pairs := strings.Split(strings.TrimSuffix(name, ext), ",")
	srcname = pairs[0] + filepath.Ext(name)
	sizePairs := strings.Split(pairs[1], "x")
	sizes = make([]int, 2)
	sizes[0], _ = strconv.Atoi(sizePairs[0])
	sizes[1], _ = strconv.Atoi(sizePairs[1])
	return
}
