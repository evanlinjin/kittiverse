package genetics

import (
	"encoding/hex"
	"errors"
)

const (
	AlleleLen = 2
)

var (
	ErrInvalidHexLen = errors.New("invalid hex length")
)

type Allele [AlleleLen]byte

func (a Allele) ToUint16() uint16 {
	return uint16(256*int(a[0]) + int(a[1]))
}

func (a Allele) ToHex() string {
	return hex.EncodeToString(a[:])
}

func (a *Allele) FromUint16(n uint16) {
	a[0], a[1] = byte(n/256), byte(n%256)
}

func (a *Allele) FromHex(hs string) error {
	h, e := hex.DecodeString(hs)
	if e != nil {
		return e
	}
	return a.Set(h)
}

func (a *Allele) Set(b []byte) error {
	if len(b) != AlleleLen {
		return ErrInvalidHexLen
	}
	copy(a[:], b[:])
	return nil
}

func (a Allele) Increment() Allele {
	if a[0]++; a[0] == 0 {
		a[1]++
	}
	return a
}

func (a Allele) Decrement() Allele {
	if a[0]--; a[0] == 255 {
		a[1]--
	}
	return a
}
