package util

import (
	"image"
	"image/color"
	"image/png"
	"net/http"
)

// P holds a pixel, and its name is shorter than P :D
type P struct {
	X int
	Y int
}

// FaviconData holds pixels to be 'coloured in' using a particular colour.
type FaviconData struct {
	Pixels []P
}

var fav8K = FaviconData{
	Pixels: []P{
		// col1
		P{2, 4}, P{2, 5}, P{2, 9}, P{2, 10},
		// col2
		P{3, 3}, P{3, 4}, P{3, 5}, P{3, 6}, P{3, 8}, P{3, 9}, P{3, 10}, P{3, 11},
		// col3
		P{4, 3}, P{4, 7}, P{4, 8}, P{4, 11},
		// col4
		P{5, 3}, P{5, 6}, P{5, 7}, P{5, 11},
		// col4
		P{6, 4}, P{6, 5}, P{6, 9}, P{6, 10},
		// col5
		P{7, 3}, P{7, 4}, P{7, 5}, P{7, 6}, P{7, 7}, P{7, 8}, P{7, 9}, P{7, 10}, P{7, 11},
		// col6
		P{8, 3}, P{8, 4}, P{8, 5}, P{8, 6}, P{8, 7}, P{8, 8}, P{8, 9}, P{8, 10}, P{8, 11},
		// col7
		P{9, 7}, P{9, 8}, P{9, 9},
		// col8
		P{10, 4}, P{10, 5}, P{10, 6}, P{10, 8}, P{10, 9}, P{10, 10},
		// col9
		P{11, 3}, P{11, 4}, P{11, 5}, P{11, 9}, P{11, 10}, P{11, 11},
		// col10
		P{12, 10}, P{12, 11},
	},
}

var favIcon *image.RGBA

func init() {
	ul := image.Point{0, 0}
	br := image.Point{16, 16}

	img := image.NewRGBA(image.Rectangle{ul, br})

	for _, p := range fav8K.Pixels {
		img.Set(p.X, p.Y, color.White)
	}

	favIcon = img

}

// GenerateFaviconPNG Generates and writes out a PNG favicon
func GenerateFaviconPNG(w http.ResponseWriter) {
	png.Encode(w, favIcon)
}
