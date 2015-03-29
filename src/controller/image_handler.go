package controller

import (
	"blog"
	"github.com/disintegration/imaging"
	// "github.com/nfnt/resize"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"time"
)

// horse.jpg
// horse,800x600.jpg
// horse,800x0.jpg
// horse,0x600.jpg
func ImageResize(dirname string) martini.Handler {
	re := regexp.MustCompile(`(.+),([0-9]+)x([0-9]+)(\.\w+)$`)
	return func(w http.ResponseWriter, r *http.Request) {
		urlpath := r.URL.Path
		matches := re.FindStringSubmatch(path.Base(urlpath))
		if len(matches) != 5 {
			return // skip
		}

		// matches = [horse,200x500.jpg horse 200 500 jpg]
		basename := matches[1] + matches[4]
		width, _ := strconv.Atoi(matches[2])
		height, _ := strconv.Atoi(matches[3])
		filename := path.Join(blog.Path(dirname), path.Dir(urlpath), basename)
		log.Printf("resize image: %s/%s, %dx%d\n", path.Dir(urlpath), basename, width, height)

		srcimg, err := imaging.Open(filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dstimg := imaging.Resize(srcimg, width, height, imaging.Box)
		gmtloc, _ := time.LoadLocation("GMT")
		w.Header().Set("Expires", time.Now().In(gmtloc).Add(time.Hour*24*7).Format(time.RFC1123))

		if err := imaging.Encode(w, dstimg, imaging.JPEG); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
