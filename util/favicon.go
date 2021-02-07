package util

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/http"
)

// P holds a pixel, and its name is shorter than Pixel :D
type P struct {
	X int
	Y int
}

type faviconData struct {
	Pixels []P
}

var fav8K = faviconData{
	Pixels: []P{
		{2, 4}, {2, 5}, {2, 9}, {2, 10}, // col1
		{3, 3}, {3, 4}, {3, 5}, {3, 6}, {3, 8}, {3, 9}, {3, 10}, {3, 11}, // col2
		{4, 3}, {4, 7}, {4, 8}, {4, 11},
		{5, 3}, {5, 6}, {5, 7}, {5, 11},
		{6, 4}, {6, 5}, {6, 9}, {6, 10},
		{7, 3}, {7, 4}, {7, 5}, {7, 6}, {7, 7}, {7, 8}, {7, 9}, {7, 10}, {7, 11},
		{8, 3}, {8, 4}, {8, 5}, {8, 6}, {8, 7}, {8, 8}, {8, 9}, {8, 10}, {8, 11},
		{9, 7}, {9, 8}, {9, 9},
		{10, 4}, {10, 5}, {10, 6}, {10, 8}, {10, 9}, {10, 10},
		{11, 3}, {11, 4}, {11, 5}, {11, 9}, {11, 10}, {11, 11},
		{12, 10}, {12, 11}, // col11
	},
}

var favIconBytes = new(bytes.Buffer)

func init() {
	ul := image.Point{0, 0}
	br := image.Point{16, 16}

	img := image.NewRGBA(image.Rectangle{ul, br})

	for _, p := range fav8K.Pixels {
		img.Set(p.X, p.Y, color.Gray16{0x5555})
	}

	png.Encode(favIconBytes, img)

}

// WriteFaviconPNG writes favIcon PNG bytes to the supplised ResponseWriter
func WriteFaviconPNG(w http.ResponseWriter) {
	w.Write(favIconBytes.Bytes())
}
