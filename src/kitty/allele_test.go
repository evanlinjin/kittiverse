package kitty

import (
	"fmt"
	"testing"
)

func TestGenome_FromUint16(t *testing.T) {
	cases := []struct {
		in  uint16
		exp Allele
	}{
		{0, Allele{0, 0}},
		{65535, Allele{255, 255}},
		{123, Allele{0, 123}},
		{511, Allele{1, 255}},
		{500, Allele{1, 244}},
	}
	var g Allele
	for i, c := range cases {
		if g.FromUint16(c.in); g != c.exp {
			t.Error(tPrint(i, c.in, c.exp, g))
		} else {
			t.Logf(tPrint(i, c.in, c.exp, g))
		}
	}
}

func TestGenome_ToUint16(t *testing.T) {
	cases := []struct {
		in  Allele
		exp uint16
	}{
		{Allele{0, 0}, 0},
		{Allele{255, 255}, 65535},
		{Allele{0, 123}, 123},
		{Allele{1, 255}, 511},
		{Allele{1, 244}, 500},
	}
	for i, c := range cases {
		if got := c.in.ToUint16(); got != c.exp {
			t.Error(tPrint(i, c.in, c.exp, got))
		} else {
			t.Log(tPrint(i, c.in, c.exp, got))
		}
	}
}

func tPrint(i int, in, exp, got interface{}) string {
	return fmt.Sprintf(
		"[%d] in(%v) expected(%v) got(%v)",
		i, in, exp, got)
}
