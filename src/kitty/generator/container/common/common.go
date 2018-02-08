package common

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	VersionLen = 2
)

var (
	ErrInvalidVersion = errors.New("invalid version")
	ErrInvalidSize    = errors.New("invalid size")
	ErrAlreadyExists  = errors.New("already exists")
)

func EmptyHash() cipher.SHA256 {
	return cipher.SHA256{}
}
