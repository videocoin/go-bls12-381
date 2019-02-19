package bls12

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
)

const fqNumWords = 6

var (
	fq0                 = fq{0}
	fq1                 = newFq(big1)
	fqNeg1              = new(fq)
	fqSqrtNeg3          = new(fq)
	fqInv2              = new(fq)
	fqHalfSqrNeg3Minus1 = new(fq)
)

type (
	// fq is a finite field with prime number q elements
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
	return elm == fq0
}

func (elm *fq) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, *elm); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// String satisfies the Stringer interface.
func (elm *fq) String() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x", elm[5], elm[4], elm[3], elm[2], elm[1], elm[0])
}
