package common

import (
	"bytes"
	"errors"
	"github.com/kittycash/kittiverse/src/kitty/generator/container"
	"github.com/skycoin/skycoin/src/cipher"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
)

const (
	XpxLen = 1200
	YpxLen = 1200

	VersionLen = 2
)

var (
	ErrInvalidVersion = errors.New("invalid version")
	ErrInvalidSize    = errors.New("invalid size")
	ErrAlreadyExists  = errors.New("already exists")
	ErrDoesNotExist   = errors.New("does not exist")
)

func EmptyHash() cipher.SHA256 {
	return cipher.SHA256{}
}

func GetImage(ic container.Images, hash cipher.SHA256) (image.Image, error) {
	raw, ok := ic.Get(hash)
	if !ok {
		return nil, ErrDoesNotExist
	}
	return png.Decode(bytes.NewReader(raw))
}

func GetRawImageFromFile(path string) ([]byte, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	raw, e := ioutil.ReadAll(f)
	if e != nil {
		return nil, e
	}
	return raw, nil
}

func DrawArea(dst draw.Image, bg, area image.Image) {
	draw.DrawMask(dst, dst.Bounds(), bg, image.ZP, area, image.ZP, draw.Over)
}

func DrawOutline(dst draw.Image, outline image.Image) {
	draw.Draw(dst, dst.Bounds(), outline, image.ZP, draw.Over)
}

func EmptyImage() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, XpxLen, YpxLen))
}
