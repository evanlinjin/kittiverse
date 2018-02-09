package container

import "github.com/skycoin/skycoin/src/cipher"

type Images interface {
	Version() uint16
	Import(raw []byte) error
	Export() []byte
	Add(raw []byte) (cipher.SHA256, error)
	Remove(hash cipher.SHA256)
	Get(hash cipher.SHA256) ([]byte, bool)
	GetOrAdd(raw []byte) cipher.SHA256
}
