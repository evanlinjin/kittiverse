package imager

import (
	"github.com/kittycash/kittiverse/src/imager/layer"
	"github.com/skycoin/skycoin/src/cipher"
)

type LayerType struct {
	Name             string
	Attributes       []Attribute                  // Organised alphabetically.
	attributesByHash map[cipher.SHA256]*Attribute `enc:"-"`
	attributesByName map[string]*Attribute        `enc:"-"`
}

type Attribute struct {
	Name         string
	Breeds       []Breed                         // Organised alphabetically.
	BreedsByHash map[cipher.SHA256]*LayerCabinet `enc:"-"`
	BreedsByName map[string]*LayerCabinet        `enc:"-"`
}

type Breed struct {
	Hash cipher.SHA256
	Name string
}

type LayerCabinet struct {
	Breed *Breed
	Parts []*Layer
}

type Layer struct {
	Area          cipher.SHA256
	AreaPlacement *layer.Placement

	Outline          cipher.SHA256
	OutlinePlacement *layer.Placement

	PartName string // PartA, PartB, etc.
}
