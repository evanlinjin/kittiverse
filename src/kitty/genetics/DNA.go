package genetics

import "encoding/hex"

// DNAPos specifies a position in the kitty DNA.
type DNAPos int

const (
	DNAVersionPos     DNAPos = iota
	DNABreedPos              = iota*6 - 5
	DNABodyAttrPos           = iota*6 - 5
	DNABodyColorAPos         = iota*6 - 5
	DNABodyColorBPos         = iota*6 - 5
	DNABodyPatternPos        = iota*6 - 5
	DNAEarsAttrPos           = iota*6 - 5
	DNAEyesAttrPos           = iota*6 - 5
	DNAEyesColorPos          = iota*6 - 5
	DNANoseAttrPos           = iota*6 - 5
	DNATailAttrPos           = iota*6 - 5
	DNAReservedAPos          = iota*6 - 5
	DNAReservedBPos          = iota*6 - 5
	DNALen                   = iota*6 - 5
)

// DNA represents a kitty's DNA and contains the genotypes of the kitty.
// A kittycash genotype is made up of 3 alleles (not 2 like real biology).
// The right-most allele will always be the dominant allele.
//		[                (    0)] DNA version (current: 0).
//		[( 1, 2),( 3, 4),( 5, 6)] Breed.
//		[( 7, 8),( 9,10),(11,12)] Body attribute.
//		[(13,14),(15,16),(17,18)] Body color A.
//		[(19,20),(21,22),(23,24)] Body color B.
//		[(25,26),(27,28),(29,30)] Body pattern.
//		[(31,32),(33,34),(35,36)] Ears attribute.
//		[(37,38),(39,40),(41,42)] Eyes attribute.
//		[(43,44),(45,46),(47,48)] Eyes color.
//		[(49,50),(51,52),(53,54)] Nose attribute.
//		[(55,56),(57,58),(59,60)] Tail attribute.
//		[(61,62),(63,64),(65,66)] (Reserved A).
//		[(67,68),(69,70),(71,72)] (Reserved B).
type DNA [DNALen]byte

func (d DNA) ToHex() string {
	return hex.EncodeToString(d[:])
}

func (d DNA) FromHex(hs string) error {
	h, e := hex.DecodeString(hs)
	if e != nil {
		return e
	}
	return d.Set(h)
}

func (d *DNA) Set(b []byte) error {
	if len(b) != DNALen {
		return ErrInvalidHexLen
	}
	copy(d[:], b[:])
	return nil
}

func (d DNA) GetGenotype(pos DNAPos) Genotype {
	return d[pos:pos+GenotypeLen]
}

func (d DNA) GetPhenotype(pos DNAPos) (a Allele) {
	copy(a[:], d[pos+4:pos+AlleleLen])
	return
}