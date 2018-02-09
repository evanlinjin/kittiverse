package genetics

import (
	"encoding/hex"
	"encoding/json"
	"errors"
)

const (
	AlleleLen = 2
)

var (
	ErrInvalidHexLen = errors.New("invalid hex length")
)

type Allele [AlleleLen]byte

func NewAlleleFromUint16(n uint16) Allele {
	return Allele{byte(n / 256), byte(n % 256)}
}

func NewAlleleFromHex(hs string) (Allele, error) {
	var a Allele
	h, e := hex.DecodeString(hs)
	if e != nil {
		return a, e
	}
	e = a.Set(h)
	return a, e
}

func (a Allele) ToUint16() uint16 {
	return uint16(256*int(a[0]) + int(a[1]))
}

func (a Allele) ToHex() string {
	return hex.EncodeToString(a[:])
}

func (a Allele) String() string {
	return a.ToHex()
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

type AlleleRange struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

type AlleleRanges struct {
	Breed         AlleleRange `json:"breed"`
	BodyAttribute AlleleRange `json:"body_attribute"`
	BodyColorA    AlleleRange `json:"body_color_a"`
	BodyColorB    AlleleRange `json:"body_color_b"`
	BodyPattern   AlleleRange `json:"body_pattern"`
	EarsAttribute AlleleRange `json:"ears_attribute"`
	EyesAttribute AlleleRange `json:"eyes_attribute"`
	EyesColor     AlleleRange `json:"eyes_color"`
	NoseAttribute AlleleRange `json:"nose_attribute"`
	TailAttribute AlleleRange `json:"tail_attribute"`
}

func (r *AlleleRanges) String(pretty bool) string {
	if pretty {
		data, _ := json.MarshalIndent(r, "", "  ")
		return string(data)
	} else {
		data, _ := json.Marshal(r)
		return string(data)
	}
}
