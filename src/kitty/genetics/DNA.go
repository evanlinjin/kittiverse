package genetics

import "encoding/hex"

// DNAPos specifies a position in the kitty DNA.
type DNAPos int

func (p DNAPos) String() string {
	return string(dnaStringArray[p])
}

type DNAPosString string

const (
	DNAVersionPos     DNAPos = iota
	DNABreedPos       DNAPos = iota*6 - 5
	DNABodyAttrPos    DNAPos = iota*6 - 5
	DNABodyColorAPos  DNAPos = iota*6 - 5
	DNABodyColorBPos  DNAPos = iota*6 - 5
	DNABodyPatternPos DNAPos = iota*6 - 5
	DNAEarsAttrPos    DNAPos = iota*6 - 5
	DNAEyesAttrPos    DNAPos = iota*6 - 5
	DNAEyesColorPos   DNAPos = iota*6 - 5
	DNANoseAttrPos    DNAPos = iota*6 - 5
	DNATailAttrPos    DNAPos = iota*6 - 5
	DNAReservedAPos   DNAPos = iota*6 - 5
	DNAReservedBPos   DNAPos = iota*6 - 5
	DNALen            int    = iota*6 - 5
)

var dnaStringArray = [...]DNAPosString{
	DNAVersionPos:     "version",
	DNABreedPos:       "breed",
	DNABodyAttrPos:    "body",
	DNABodyColorAPos:  "bodyColorA",
	DNABodyColorBPos:  "bodyColorB",
	DNABodyPatternPos: "bodyPattern",
	DNAEarsAttrPos:    "ears",
	DNAEyesAttrPos:    "eyes",
	DNAEyesColorPos:   "eyesColor",
	DNANoseAttrPos:    "nose",
	DNATailAttrPos:    "tail",
	DNAReservedAPos:   "",
	DNAReservedBPos:   "",
}

var dnaPosArray = [...]DNAPos{
	DNABreedPos,
	DNABodyAttrPos,
	DNABodyColorAPos,
	DNABodyColorBPos,
	DNABodyPatternPos,
	DNAEarsAttrPos,
	DNAEyesAttrPos,
	DNAEyesColorPos,
	DNANoseAttrPos,
	DNATailAttrPos,
}

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

func NewDNAFromHex(hs string) (DNA, error) {
	var dna DNA
	h, e := hex.DecodeString(hs)
	if e != nil {
		return dna, e
	}
	if e := dna.Set(h); e != nil {
		return dna, e
	}
	return dna, nil
}

func (d DNA) Hex() string {
	return hex.EncodeToString(d[:])
}

func (d *DNA) Set(b []byte) error {
	if len(b) != DNALen {
		return ErrInvalidHexLen
	}
	copy(d[:], b[:])
	return nil
}

func (d *DNA) SetVersion(v byte) {
	d[DNAVersionPos] = v
}

func (d *DNA) SetGenotype(pos DNAPos, a2, a1, a0 Allele) {
	copy(d[pos+0:pos+2], a2[:])
	copy(d[pos+2:pos+4], a1[:])
	copy(d[pos+4:pos+6], a0[:])
}

func (d *DNA) SetRandomGenotype(pos DNAPos, ar AlleleRange) {
	d.SetGenotype(pos, ar.GetRandom(), ar.GetRandom(), ar.GetRandom())
}

func (d DNA) GetGenotype(pos DNAPos) Genotype {
	return d[pos : pos+GenotypeLen]
}

func (d DNA) GetPhenotype(pos DNAPos) (a Allele) {
	copy(a[:], d[pos+4:pos+4+AlleleLen])
	return
}

type BreakdownSub struct{
	Version string `json:"version"`
	Breed   *GenotypeBreakdown `json:"breed"`
	BodyAttribute  *GenotypeBreakdown `json:"body_attribute"`
	BodyColorA     *GenotypeBreakdown `json:"body_color_a"`
	BodyColorB     *GenotypeBreakdown `json:"body_color_b"`
	BodyPattern    *GenotypeBreakdown `json:"body_pattern"`
	EarsAttribute  *GenotypeBreakdown `json:"ears_attribute"`
	EyesAttribute  *GenotypeBreakdown `json:"eyes_attribute"`
	EyesColor      *GenotypeBreakdown `json:"eyes_color"`
	NoseAttribute  *GenotypeBreakdown `json:"nose_attribute"`
	TailAttribute  *GenotypeBreakdown `json:"tail_attribute"`
	ReservedA      *GenotypeBreakdown `json:"reserved_a"`
	ReservedB      *GenotypeBreakdown `json:"reserved_b"`
}

type DNABreakdown struct {
	Hex string `json:"hex"`
	Breakdown BreakdownSub `json:"breakdown"`
}

func (d DNA) Breakdown() *DNABreakdown {
	return &DNABreakdown{
		Hex: d.Hex(),
		Breakdown: BreakdownSub{
			Version: hex.EncodeToString(d.GetGenotype(DNAVersionPos)[:1]),
			Breed: d.GetGenotype(DNABreedPos).Breakdown(),
			BodyAttribute: d.GetGenotype(DNABodyAttrPos).Breakdown(),
			BodyColorA: d.GetGenotype(DNABodyColorAPos).Breakdown(),
			BodyColorB: d.GetGenotype(DNABodyColorBPos).Breakdown(),
			BodyPattern: d.GetGenotype(DNABodyPatternPos).Breakdown(),
			EarsAttribute: d.GetGenotype(DNAEarsAttrPos).Breakdown(),
			EyesAttribute: d.GetGenotype(DNAEyesAttrPos).Breakdown(),
			EyesColor: d.GetGenotype(DNAEyesColorPos).Breakdown(),
			NoseAttribute: d.GetGenotype(DNANoseAttrPos).Breakdown(),
			TailAttribute: d.GetGenotype(DNATailAttrPos).Breakdown(),
			ReservedA: d.GetGenotype(DNAReservedAPos).Breakdown(),
			ReservedB: d.GetGenotype(DNAReservedBPos).Breakdown(),
		},
	}
}