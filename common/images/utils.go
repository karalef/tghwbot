package images

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"
)

// Resize image.
func Resize(src image.Image, width, height int) image.Image {
	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(resized, resized.Rect, src, src.Bounds(), draw.Over, nil)
	return resized
}

// Crop image.
func Crop(img image.Image, rect image.Rectangle) image.Image {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	simg, ok := img.(subImager)
	if !ok {
		return CropCopy(img, rect)
	}
	return simg.SubImage(rect)
}

// CropCopy creates copy of image part.
func CropCopy(img image.Image, rect image.Rectangle) image.Image {
	rgbaImg := &image.RGBA{}
	if rect = rect.Intersect(img.Bounds()); !rect.Empty() {
		rgbaImg = image.NewRGBA(rect)
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				rgbaImg.Set(x, y, img.At(x, y))
			}
		}
	}
	return rgbaImg
}

// DrawClipCircle clips src image with circle and draws it.
func DrawClipCircle(dst draw.Image, offset image.Point, src image.Image, center image.Point, r int) {
	mask := &circle{
		p: center,
		r: r,
	}
	draw.DrawMask(dst, dst.Bounds().Add(offset), src, image.ZP, mask, image.ZP, draw.Over)
}

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
