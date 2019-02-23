package bls12

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

const (
	// fqLen is the expected length of a field element
	fqLen = 6
)

var errOutOfBounds = errors.New("value is not an element of the finite field of order q")

var (
	fq0    = fq{0}
	fq1, _ = fqFromBig(big1) // TODO This is wrong (montgomery form)
	fqR2   = fq{}            // TODO
)

// fq is an element of the finite field of order q.
type (
	fq      [fqLen]uint64
	fqLarge [fqLen * 2]uint64
)

// equal checks if the field elements are equal.
func (fq fq) equal(b fq) bool {
	for i, fqi := range fq {
		if fqi != b[i] {
			return false
		}
	}

	return true
}

// Hex returns the field element in the hexadecimal base
func (fq *fq) Hex() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x", fq[5], fq[4], fq[3], fq[2], fq[1], fq[0])
}

// String satisfies the Stringer interface.
func (fq *fq) String() string {
	return fq.Hex()
}

// isFieldElement checks if value is within the field bounds.
func isFieldElement(value *big.Int) bool {
	return (value.Sign() >= 0) && (value.Cmp(bigQ) < 0)
}

// fqFromBig converts a big integer to a field element.
func fqFromBig(value *big.Int) (fq, error) {
	if !isFieldElement(value) {
		return fq{}, errOutOfBounds
	}

	fq := fq{0}
	words := value.Bits()
	numWords := len(words)
	if strconv.IntSize == 64 {
		for i := 0; i < numWords; i++ {
			fq[i] = uint64(words[i])
		}
	} else {
		for i := 0; i < numWords; i++ {
			fq[i/2] |= uint64(words[i]) << uint(32*(i%2))
		}
	}

	// TODO - Montgomery form
	// montEncode(fq)

	return fq, nil
}

/*
func montEncode(x *fq) {
	fqMul(x, x, r2)
}

func montDecode(x *fq) {

}
*/
