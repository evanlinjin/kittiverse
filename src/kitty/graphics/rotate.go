package layer

import (
	"github.com/kittycash/kittiverse/src/kitty/graphics/graphics-go/graphics"
	"image"
	"image/draw"
	"math"
)

// Rotate rotates an image while keeping all the pixels of the original image.
func Rotate(src image.Image, radians float64) (draw.Image, error) {

	// Fixes.

	radians = math.Mod(radians, math.Pi*2)

	// Find destination size.

	srcBounds := src.Bounds()
	srcX := float64(getRectWidth(srcBounds)) / 2
	srcY := float64(getRectHeight(srcBounds)) / 2

	z := math.Sqrt(math.Pow(srcX, 2) + math.Pow(srcY, 2))

	var angle float64
	if radians >= 0 && radians < math.Pi/2 || radians >= math.Pi && radians < math.Pi*3/2 {
		angle = radians
	} else {
		angle = -radians
	}

	dstY := z * math.Sin(math.Atan(srcY/srcX)+angle)
	dstX := z * math.Sin(math.Atan(srcX/srcY)+angle)

	// Create destination image.

	dst := image.NewRGBA(image.Rect(0, 0, int(math.Ceil(dstX*2)), int(math.Ceil(dstY*2))))
	e := graphics.Rotate(dst, src, &graphics.RotateOptions{Angle: radians})
	return dst, e
}
