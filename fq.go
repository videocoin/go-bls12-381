package bls12

import (
	"fmt"
	"math/big"
	"strconv"
)

const fqNumWords = 6

var fqZero = fq{0}

type (
	fq      [fqNumWords]uint64
	fqLarge [2 * fqNumWords]uint64
)

func newFq(n *big.Int) fq {
	fq := fq{}
	words := n.Bits()
	numWords := len(words)
	if strconv.IntSize == 64 {
		for i := 0; i < numWords && i < fqNumWords; i++ {
			fq[i] = uint64(words[i])
		}
	} else {
		for i := 0; i < numWords && i < fqNumWords*2; i++ {
			fq[i/2] = uint64(words[i]) << uint(32*(i%2))
		}
	}
	return fq
}

func (elm fq) isZero() bool {
	return elm == fqZero
}

// String satisfies the Stringer interface.
func (elm *fq) String() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x", elm[5], elm[4], elm[3], elm[2], elm[1], elm[0])
}
