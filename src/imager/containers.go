package imager

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"io"
	"io/ioutil"
)

/*
	<<< IMAGE CONTAINER >>>
*/

// ImageContainer stores raw images via hash key.
// TODO (evanlinjin): Make this more memory-efficient (use 'github.com/boltdb/bolt' ?)
type ImageContainer struct {
	dict map[cipher.SHA256][]byte
}

// FromReader loads ImageContainer from reader.
func (c *ImageContainer) FromReader(r io.Reader) error {
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

// ToWriter writes ImageContainer to writer.
func (c *ImageContainer) ToWriter(w io.Writer) error {
	out, i := make([][]byte, len(c.dict)), 0
	for _, v := range c.dict {
		out[i], i = v, i+1
	}
	_, e := w.Write(encoder.Serialize(out))
	return e
}

// Store stores an image in the internal image container.
func (c *ImageContainer) Store(raw []byte) cipher.SHA256 {
	hash := cipher.SumSHA256(raw)
	c.dict[hash] = raw
	return hash
}

// Delete removes an image from internal image container.
func (c *ImageContainer) Delete(hash cipher.SHA256) {
	delete(c.dict, hash)
}

// Load obtains an image from internal image container.
func (c *ImageContainer) Load(hash cipher.SHA256) ([]byte, bool) {
	v, ok := c.dict[hash]
	if !ok {
		return nil, false
	}
	return v, true
}

/*
	<<< LAYER TYPES CONTAINER >>>
*/

// LayerTypesContainer contains all the layer_types, nicely indexed and sorted.
type LayerTypesContainer struct {
	list []LayerType
	dict map[string]*LayerType
}

// FromReader loads LayerTypesContainer from reader.
func (c *LayerTypesContainer) FromReader(r io.Reader) error {
	raw, e := ioutil.ReadAll(r)
	if e != nil {
		return e
	}
	if e := encoder.DeserializeRaw(raw, &c.list); e != nil {
		return e
	}
	c.dict = make(map[string]*LayerType)
	for i := 0; i < len(c.list); i++ {
		p := &c.list[0]
		c.dict[p.Name] = p
	}
	return nil
}

// ToWriter writes LayerTypesContainer to writer.
func (c *LayerTypesContainer) ToWriter(w io.Writer) error {
	_, e := w.Write(encoder.Serialize(c.list))
	return e
}

// Store stores a LayerType in the container.
func (c *LayerTypesContainer) Store(lt LayerType) error {
	return nil
}

// Load loads  from the container.
func (c *LayerTypesContainer) Load(name string) (*LayerType, bool) {
	v, ok := c.dict[name]
	return v, ok
}
