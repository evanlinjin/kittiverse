package container

import (
	"github.com/kittycash/kittiverse/src/kitty/genetics"
	"image"
)

type Layers interface {
	Version() uint16
	Import(raw []byte) error
	Export() []byte
	Compile(rootDir string, images Images) error
	GetAlleleRanges() *genetics.AlleleRanges
	GenerateKitty(images Images, dna genetics.DNA) (image.Image, error)
}
