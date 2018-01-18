package incubator

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type KittyGenerator struct {
	c    *KittyGenSpecs
	skin image.Image
	l    []image.Image
}

func NewKittyGenerator(c *KittyGenSpecs) (*KittyGenerator, error) {
	if e := c.Init(); e != nil {
		return nil, e
	}

	skinPaths := c.GetSkinPaths()

	skin := image.NewRGBA(image.Rect(ImageMin, ImageMin, ImageMax, ImageMax))

	fill, e := getImageFromPath(skinPaths.Color)
	if e != nil {
		return nil, e
	}
	draw.Draw(skin, skin.Bounds(), fill, fill.Bounds().Min, draw.Over)

	if skinPaths.HasPattern() {
		pattern, e := getImageFromPath(skinPaths.Pattern)
		if e != nil {
			return nil, e
		}
		draw.Draw(skin, skin.Bounds(), pattern, pattern.Bounds().Min, draw.Over)
	}

	return &KittyGenerator{
		c:    c,
		skin: skin,
	}, nil
}

func (g *KittyGenerator) addLayer(part KittyPart) error {
	var (
		partSpecs = part.Specs()
		partPaths = g.c.GetPartPath(partSpecs)
	)
	if partPaths == nil {
		return nil
	}

	// Outline.
	outlineImg, e := getImageFromPath(partPaths.Outline)
	if e != nil {
		outlineImg, e = getImageFromPath(partPaths.OutlineAlt)
		if e != nil {
			return e
		}
	}

	// Area.
	areaImg, e := getImageFromPath(partPaths.Area)
	if e != nil {
		// No area.
		g.l = append(g.l, outlineImg)
		return nil
	}

	var bg = g.skin.(*image.RGBA)
	if partSpecs.IsAccessory() && partPaths.HasColor() {
		bg = image.NewRGBA(image.Rect(ImageMin, ImageMin, ImageMax, ImageMax))
		color, e := getImageFromPath(partPaths.Color)
		if e != nil {
			return e
		}
		draw.Draw(bg, bg.Bounds(), color, image.ZP, draw.Over)
	}

	layer := image.NewRGBA(image.Rect(ImageMin, ImageMin, ImageMax, ImageMax))
	draw.DrawMask(layer, layer.Bounds(), bg, image.ZP, areaImg, image.ZP, draw.Over)
	draw.Draw(layer, layer.Bounds(), outlineImg, image.ZP, draw.Over)
	g.l = append(g.l, layer)

	return nil
}

func (g *KittyGenerator) Compile() (image.Image, error) {

	e := RangeKittyParts(func(part KittyPart) error {
		return g.addLayer(part)
	})
	if e != nil {
		return nil, e
	}

	out := image.NewRGBA(image.Rect(ImageMin, ImageMin, ImageMax, ImageMax))

	for i, img := range g.l {
		fmt.Println(i, img.Bounds())
		draw.Draw(out, out.Bounds(), img, img.Bounds().Min, draw.Over)
	}

	return out, nil
}

func (g *KittyGenerator) CompileToFile(name string) (image.Image, error) {
	img, e := g.Compile()
	if e != nil {
		return nil, e
	}
	f, e := os.Create(name)
	if e != nil {
		return nil, e
	}
	return img, png.Encode(f, img)
}

func GenerateKitty(specs *KittyGenSpecs, toFile bool, fName string) (image.Image, error) {
	gen, e := NewKittyGenerator(specs)
	if e != nil {
		return nil, e
	}
	if toFile {
		return gen.CompileToFile(fName)
	} else {
		return gen.Compile()
	}
}

/*
	<< HELPER FUNCTIONS >>
*/

func getImageFromPath(path string) (image.Image, error) {
	if f, e := os.Open(path); e != nil {
		return nil, e
	} else if img, e := png.Decode(f); e != nil {
		return nil, e
	} else if s := img.Bounds().Size(); s.X != ImageMax || s.Y != ImageMax {
		return nil, errors.New("skin fill image has invalid dimensions")
	} else {
		return img, nil
	}
}
