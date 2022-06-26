package fonts

import (
	_ "embed" // for fonts

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

//go:embed Roboto-Regular.ttf
var robotoTtf []byte

//go:embed Roboto-Bold.ttf
var robotoBoldTtf []byte

var goFont, _ = truetype.Parse(goregular.TTF)
var roboto, _ = truetype.Parse(robotoTtf)
var robotoBold, _ = truetype.Parse(robotoBoldTtf)

// GoFontFace returns goregular font face.
func GoFontFace(size float64) font.Face {
	return truetype.NewFace(goFont, &truetype.Options{Size: size})
}

// RobotoFace returns Roboto font face.
func RobotoFace(size float64) font.Face {
	return truetype.NewFace(roboto, &truetype.Options{Size: size})
}

// RobotoBoldFace returns Roboto Bold font face.
func RobotoBoldFace(size float64) font.Face {
	return truetype.NewFace(robotoBold, &truetype.Options{Size: size})
}

// Height returns font height.
func Height(ff font.Face) int {
	return ff.Metrics().Height.Ceil()
}

// StringWidth returns rendering width of provided string.
func StringWidth(ff font.Face, s string) int {
	w := 0
	for _, r := range s {
		_, a, _ := ff.GlyphBounds(r)
		w += a.Round()
	}
	return w
}
