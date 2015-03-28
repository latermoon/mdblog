package controller

import (
	"blog"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func imageResizeHandler(w http.ResponseWriter, r *http.Request) {
	dir, _, srcname, sizes := fileInfo(r.URL.Path)
	var srcfile string
	if strings.HasPrefix(r.URL.Path, "/private/") {
		srcfile = filepath.Join(blog.Workspace(), dir, srcname)
	} else {
		srcfile = filepath.Join(blog.Workspace(), "public", dir, srcname)
	}

	log.Println(r.URL.Path, srcname, sizes)
	mimg, err := loadImage(srcfile)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	gmtLoc, _ := time.LoadLocation("GMT")
	w.Header().Set("Expires", time.Now().In(gmtLoc).Add(time.Hour*24*7).Format(time.RFC1123))
	rimg := resize.Thumbnail(uint(sizes[0]), uint(sizes[1]), mimg, resize.Lanczos3)
	encodeImage(w, rimg, filepath.Ext(srcname))
}

func encodeImage(w io.Writer, m image.Image, ext string) {
	switch strings.ToLower(ext) {
	case ".png":
		png.Encode(w, m)
	default: // ".jpg", ".jpeg":
		jpeg.Encode(w, m, &jpeg.Options{Quality: 85})
	}
}

func loadImage(filename string) (img image.Image, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ext := filepath.Ext(filename)
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
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
