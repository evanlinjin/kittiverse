package genetics

import "encoding/hex"

const (
	GenotypeLen = 6
)

type Genotype []byte

func (g Genotype) Hex() string {
	return hex.EncodeToString(g)
}

func (g Genotype) Breakdown() *GenotypeBreakdown {
	return &GenotypeBreakdown{
		Recessive1: Allele{g[0], g[1]}.Hex(),
		Recessive2: Allele{g[2], g[3]}.Hex(),
		Dominant:   Allele{g[4], g[5]}.Hex(),
	}
}

type GenotypeBreakdown struct {
	Recessive1 string `json:"r1"`
	Recessive2 string `json:"r2"`
	Dominant   string `json:"d"`
}