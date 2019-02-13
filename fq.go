package bls12

import (
	"fmt"
	"math/big"
	"strconv"
)

const fqLen = 6

type fq [fqLen]uint64

func newFq(n *big.Int) fq {
	fq := fq{}
	words := n.Bits()
	numWords := len(words)
	if strconv.IntSize == 64 {
		for i := 0; i < numWords && i < fqLen; i++ {
			fq[i] = uint64(words[i])
		}
	} else {
		for i := 0; i < numWords && i < fqLen*2; i++ {
			fq[i/2] = uint64(words[i]) << uint(32*(i%2))
		}
	}
	return fq
}

func (elm *fq) isZero() bool {
	// TODO
	return false
}

func (elm *fq) isOne() bool {
	// TODO
	return false
}

// String satisfies the Stringer interface.
func (elm *fq) String() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x", elm[5], elm[4], elm[3], elm[2], elm[1], elm[0])
}

func fqAdd(c, a, b *fq) {}
func fqMul(c, a, b *fq) {}
func fqSub(c, a, b *fq) {}
func fqSqr(c, a *fq)    {}
func fqDbl(c, a *fq)    {}
