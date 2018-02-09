package v0

import (
	"github.com/kittycash/kittiverse/src/kitty/generator/container/common"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"gopkg.in/sirupsen/logrus.v1"
)

var log = logrus.New()

func init() {
	log.SetLevel(logrus.DebugLevel)
}

const (
	version uint16 = 0
)

type Images struct {
	Images       [][]byte
	imagesByHash map[cipher.SHA256]*[]byte `enc:"-"`
}

func NewImagesContainer() *Images {
	return &Images{
		imagesByHash: make(map[cipher.SHA256]*[]byte),
	}
}

func (ic *Images) Version() uint16 {
	return version
}

func (ic *Images) Import(raw []byte) error {
	// Check raw size.
	if len(raw) < common.VersionLen {
		return common.ErrInvalidSize
	}
	// Check version.
	var ver uint16
	if e := encoder.DeserializeRaw(raw[:common.VersionLen], &ver); e != nil {
		return e
	}
	if ver != version {
		return common.ErrInvalidVersion
	}
	// Load data.
	if e := encoder.DeserializeRaw(raw[common.VersionLen:], ic); e != nil {
		return e
	}
	// Prepare map.
	ic.imagesByHash = make(map[cipher.SHA256]*[]byte)
	for i, v := range ic.Images {
		ic.imagesByHash[cipher.SumSHA256(v)] = &ic.Images[i]
	}
	return nil
}

func (ic *Images) Export() []byte {
	return append(encoder.Serialize(version), encoder.Serialize(ic)...)
}

func (ic *Images) Add(raw []byte) (cipher.SHA256, error) {
	hash := cipher.SumSHA256(raw)
	// Check if map is prepared.
	if ic.imagesByHash == nil {
		ic.imagesByHash = make(map[cipher.SHA256]*[]byte)
	}
	// Check if already exists.
	if _, has := ic.imagesByHash[hash]; has {
		return common.EmptyHash(), common.ErrAlreadyExists
	}
	// Append.
	ic.Images = append(ic.Images, raw)
	ic.imagesByHash[hash] = &ic.Images[len(ic.Images)-1]
	return hash, nil
}

func (ic *Images) Remove(hash cipher.SHA256) {
	if _, has := ic.imagesByHash[hash]; has {
		delete(ic.imagesByHash, hash)
		for i, v := range ic.Images {
			if len(v) > 0 && cipher.SumSHA256(v) == hash {
				ic.Images[i] = []byte{}
				return
			}
		}
	}
}

func (ic *Images) Get(hash cipher.SHA256) ([]byte, bool) {
	v, ok := ic.imagesByHash[hash]
	if !ok {
		return nil, false
	}
	return *v, true
}

func (ic *Images) GetOrAdd(raw []byte) cipher.SHA256 {
	hash := cipher.SumSHA256(raw)
	// Check if map is prepared.
	if ic.imagesByHash == nil {
		ic.imagesByHash = make(map[cipher.SHA256]*[]byte)
	}
	// Check if already exists.
	if _, has := ic.imagesByHash[hash]; has {
		return hash
	}
	// Append.
	ic.Images = append(ic.Images, raw)
	ic.imagesByHash[hash] = &ic.Images[len(ic.Images)-1]
	return hash
}
