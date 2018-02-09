package generator

import (
	"github.com/kittycash/kittiverse/src/kitty/generator/container"
	"github.com/kittycash/kittiverse/src/kitty/genetics"
	"github.com/sirupsen/logrus"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"image"
	"io"
)

type InstanceFile struct {
	Images []byte
	Layers []byte
}

type Instance struct {
	log *logrus.Logger   // logging.
	ic  container.Images // contains all images.
	lc  container.Layers // contains layers.
}

func NewInstance(imagesContainer container.Images, layersContainer container.Layers) *Instance {
	return &Instance{
		log: logrus.New(),
		ic:  imagesContainer,
		lc:  layersContainer,
	}
}

func (i *Instance) Import(r io.ReadCloser, size int) error {
	file := new(InstanceFile)
	if e := encoder.Deserialize(r, size, file); e != nil {
		return e
	}
	if e := i.ic.Import(file.Images); e != nil {
		return e
	}
	if e := i.lc.Import(file.Layers); e != nil {
		return e
	}
	return nil
}

func (i *Instance) Export(w io.Writer) error {
	_, e := w.Write(encoder.Serialize(InstanceFile{
		Images: i.ic.Export(),
		Layers: i.lc.Export(),
	}))
	return e
}

func (i *Instance) Compile(dir string) error {
	return i.lc.Compile(dir, i.ic)
}

func (i *Instance) GetAlleleRanges() genetics.AlleleRanges {
	return i.lc.GetAlleleRanges()
}

func (i *Instance) GenerateKitty(dna genetics.DNA) (image.Image, error) {
	return i.lc.GenerateKitty(dna)
}
