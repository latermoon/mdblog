package controller

import (
	"blog"
	"crypto/md5"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var imgRegx = regexp.MustCompile(`(.+),([0-9]{0,4})x([0-9]{0,4})(\.\w+)$`)

// horse.jpg
// horse,800x600.jpg
// horse,800x0.jpg
// horse,0x600.jpg
func ImageResize(dirname string) martini.Handler {
	wlock := sync.Mutex{}
	return func(w http.ResponseWriter, r *http.Request) {
		urlpath := r.URL.Path
		basename, width, height := splitImageUrlPath(urlpath)
		if basename == "" {
			return // skip
		}
		filename := path.Join(blog.Path(dirname), path.Dir(urlpath), basename)

		// check cache
		cachefile := imageCacheName(urlpath)

		if !fileExist(cachefile) {
			// lock write file
			wlock.Lock()
			defer wlock.Unlock()

			file, err := os.Create(cachefile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			if err := resizeImage(file, filename, width, height); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			w.Header().Set("hit", path.Base(cachefile))
		}

		w.Header().Set("Expires", expiresHeader(time.Now().Add(time.Hour*24*7)))
		http.ServeFile(w, r, cachefile)
	}
}

func resizeImage(w io.Writer, src string, width, height int) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	log.Printf("resize image: %s, %dx%d\n", path.Base(src), width, height)
	dstimg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	switch strings.ToLower(path.Ext(src)) {
	case ".png":
		return png.Encode(w, dstimg)
	default:
		return jpeg.Encode(w, dstimg, &jpeg.Options{jpeg.DefaultQuality})
	}
}

// split "horse,200x500.jpg" into  [horse.jpg 200 500]
func splitImageUrlPath(urlpath string) (basename string, width, height int) {
	// matches = [horse,200x500.jpg horse 200 500 .jpg]
	matches := imgRegx.FindStringSubmatch(path.Base(urlpath))
	if len(matches) != 5 {
		return
	}
	basename = matches[1] + matches[4]
	width, _ = strconv.Atoi(matches[2])
	height, _ = strconv.Atoi(matches[3])
	return
}

func imageCacheName(urlpath string) string {
	name := fmt.Sprintf("%x", md5.Sum([]byte(blog.Path(urlpath))))
	dir := "/tmp/" + fmt.Sprintf("%x", md5.Sum([]byte("mdblog")))
	if !fileExist(dir) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			log.Println(err)
		}
	}
	return path.Join(dir, name)
}
