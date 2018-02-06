package generator

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type LayersOfType struct {
	OfType      string
	Layers      []Layer
	layersByKey map[attributeKey]*Layer `enc:"-"` // aka, by attribute and breed
}

func (lt *LayersOfType) Init() {
	if lt.layersByKey == nil {
		lt.layersByKey = make(map[attributeKey]*Layer)
	}
	for i := 0; i < len(lt.Layers); i++ {
		layer := &lt.Layers[i]
		lt.layersByKey[layer.key()] = layer
	}
}

func (lt *LayersOfType) Add(layer Layer) error {
	if _, has := lt.layersByKey[layer.key()]; has {
		return ErrAlreadyExists
	}
	lt.Layers = append(lt.Layers, layer)
	lt.layersByKey[layer.key()] = &lt.Layers[len(lt.Layers)-1]
	return nil
}

// Layer represents a kitty layer.
// Field "Parts" represents a slice of image hashes in pairs, in which;
// 		1. each slice element represents a "part", and
//		2. within each part, the hash pair consists of;
//			[0] representing the layer "area".
//			[1] representing the layer "outline".
type Layer struct {
	OfAttribute string
	OfBreed     string
	Parts       [][2]cipher.SHA256
}

type attributeKey string

func (a *Layer) key() attributeKey {
	return attributeKey(a.OfAttribute + "_" + a.OfBreed)
}

