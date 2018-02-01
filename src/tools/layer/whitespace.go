package layer

import (
	"github.com/kittycash/kittiverse/src/incubator/types"
	"image"
	"image/draw"
)

func RemoveWhitespace(src image.Image) (image.Image, *types.LayerPlacement) {
	var (
		b    = src.Bounds()
		minX int
		maxX int
		minY int
		maxY int
	)
	for x := b.Min.X; x <= b.Max.X; x++ {
		for y := b.Min.Y; y <= b.Max.Y; y++ {
			if _, _, _, a := src.At(x, y).RGBA(); a > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	var dst = image.NewRGBA(image.Rect(minX, minY, maxX, maxY))
	draw.Draw(dst, dst.Bounds(), src, image.Pt(minX, minY), draw.Src)

	return dst, &types.LayerPlacement{
		DisplaceX: uint64((minX+maxX)/2),
		DisplaceY: uint64((minY+maxY)/2),
		ScaleX:    1,
		ScaleY:    1,
		Rotate:    0,
	}
}

func IncludeWhitespace(src image.Image, bounds image.Rectangle, at *types.LayerPlacement) image.Image {
	var (
		dst = image.NewRGBA(bounds)
		srcBound = src.Bounds()
		spX = int(at.DisplaceX) - getRectWidth(srcBound)
		spY = int(at.DisplaceY) - getRectHeight(srcBound)
	)

	Rotate(src, at.Rotate)

	draw.Draw(dst, bounds, src, image.Pt(spX, spY), draw.Src)
	//graphics.Scale()

	//graphics.Rotate()
	//graphics.Scale()

	return dst
}

func getRectWidth(rectangle image.Rectangle) int {
	return rectangle.Max.X - rectangle.Min.X
}

func getRectHeight(rectangle image.Rectangle) int {
	return rectangle.Max.Y - rectangle.Min.Y
}