package generator

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"io"
	"io/ioutil"
	"errors"
)

var (
	ErrAlreadyExists = errors.New("already exists")
)

/*
	<<< IMAGE CONTAINER >>>
*/

// ImageContainer stores raw images via hash key.
// TODO (evanlinjin): Make this more memory-efficient (use 'github.com/boltdb/bolt' ?)
type ImageContainer struct {
	dict map[cipher.SHA256][]byte
}

// Import loads ImageContainer from reader.
func (c *ImageContainer) Import(r io.Reader) error {
	raw, e := ioutil.ReadAll(r)
	if e != nil {
		return e
	}
	var out [][]byte
	if e := encoder.DeserializeRaw(raw, &out); e != nil {
		return e
	}
	c.dict = make(map[cipher.SHA256][]byte)
	for _, v := range out {
		c.dict[cipher.SumSHA256(v)] = v
	}
	return nil
}

// Export writes ImageContainer to writer.
func (c *ImageContainer) Export(w io.Writer) error {
	out, i := make([][]byte, len(c.dict)), 0
	for _, v := range c.dict {
		out[i], i = v, i+1
	}
	_, e := w.Write(encoder.Serialize(out))
	return e
}

// Add stores an image in the internal image container.
func (c *ImageContainer) Add(raw []byte) cipher.SHA256 {
	hash := cipher.SumSHA256(raw)
	c.dict[hash] = raw
	return hash
}

// Remove removes an image from internal image container.
func (c *ImageContainer) Remove(hash cipher.SHA256) {
	delete(c.dict, hash)
}

// Get obtains an image from internal image container.
func (c *ImageContainer) Get(hash cipher.SHA256) ([]byte, bool) {
	v, ok := c.dict[hash]
	if !ok {
		return nil, false
	}
	return v, true
}

/*
	<<< LAYERS CONTAINER >>>
*/

// LayersContainer contains all the layer_types, nicely indexed and sorted.
type LayersContainer struct {
	LayerTypes       []LayersOfType
	Breeds           []string
	layerTypesByName map[string]*LayersOfType `enc:"-"`
}

// Import loads LayersContainer from reader.
func (c *LayersContainer) Import(r io.Reader) error {
	raw, e := ioutil.ReadAll(r)
	if e != nil {
		return e
	}
	if e := encoder.DeserializeRaw(raw, c); e != nil {
		return e
	}
	c.layerTypesByName = make(map[string]*LayersOfType)
	for i := 0; i < len(c.LayerTypes); i++ {
		p := &c.LayerTypes[i]
		c.layerTypesByName[p.OfType] = p
	}
	return nil
}

// Export writes LayersContainer to writer.
func (c *LayersContainer) Export(w io.Writer) error {
	_, e := w.Write(encoder.Serialize(c))
	return e
}

// Add stores a LayersOfType in the container.
func (c *LayersContainer) Add(lt LayersOfType) error {
	if _, has := c.layerTypesByName[lt.OfType]; has {
		return ErrAlreadyExists
	}
	return nil
}

// Get loads from the container.
func (c *LayersContainer) Get(name string) (*LayersOfType, bool) {
	v, ok := c.layerTypesByName[name]
	return v, ok
}
