package layer

import (
	"github.com/kittycash/kittiverse/src/imager/layer/graphics-go/graphics"
	"image"
	"image/draw"
)

func Scale(src image.Image, scaleX, scaleY float64) (draw.Image, error) {
	srcBounds := src.Bounds()
	srcWidth := getRectWidth(srcBounds)
	srcHeight := getRectHeight(srcBounds)

	dstWidth := int(float64(srcWidth) * scaleX)
	dstHeight := int(float64(srcHeight) * scaleY)
	dstBounds := image.Rect(0, 0, dstWidth, dstHeight)
	dst := image.NewRGBA(dstBounds)

	if e := graphics.Scale(dst, src); e != nil {
		return nil, e
	}
	return dst, nil
}
