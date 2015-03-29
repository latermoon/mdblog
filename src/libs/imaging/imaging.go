package imaging

import (
	"image"
	// "image/gif"
	"image/jpeg"
	// "image/png"
	"io"
	"os"
)

func Open(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return image.Decode(file)
}

func Encode(w io.Writer, img image.Image) {
	jpeg.Encode(w, img, nil)
}
