package layer

import (
	"image"
	"image/draw"
)

func RemoveWhitespace(src image.Image) (image.Image, *Placement) {
	var (
		b    = src.Bounds()
		minX = -1
		minY = -1
		maxX = -1
		maxY = -1
	)
	for x := b.Min.X; x <= b.Max.X; x++ {
		for y := b.Min.Y; y <= b.Max.Y; y++ {
			if _, _, _, a := src.At(x, y).RGBA(); a > 0 {
				if x < minX || minX == -1 {
					minX = x
				}
				if y < minY || minY == -1 {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	var dst = image.NewRGBA(image.Rect(minX, minY, maxX, maxY))
	draw.Draw(dst, dst.Bounds(), src, image.Pt(minX, minY), draw.Over)

	return dst, &Placement{
		CoordX: uint64((minX+maxX)/2),
		CoordY: uint64((minY+maxY)/2),
		ScaleX: 1,
		ScaleY: 1,
		Rotate: 0,
	}
}

func IncludeWhitespace(src image.Image, bounds image.Rectangle, at *Placement) (image.Image, error) {
	var (
		dst = image.NewRGBA(bounds)
		e error
	)

	if src, e = Scale(src, at.ScaleX, at.ScaleY); e != nil {
		return nil, e
	}

	if src, e = Rotate(src, at.Rotate); e != nil {
		return nil, e
	}

	var (
		srcBound = src.Bounds()
		spX = int(at.CoordX) - getRectWidth(srcBound)/2
		spY = int(at.CoordY) - getRectHeight(srcBound)/2
	)
	draw.Draw(dst, bounds, src, image.Pt(-spX, -spY), draw.Src)

	return dst, nil
}

func getRectWidth(rectangle image.Rectangle) int {
	return rectangle.Max.X - rectangle.Min.X
}

func getRectHeight(rectangle image.Rectangle) int {
	return rectangle.Max.Y - rectangle.Min.Y
}