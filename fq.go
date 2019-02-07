package bls12

import (
	"fmt"
	"math/big"
	"strconv"
)

const (
	fqLen = 6

	decimalBase = 10
)

func bigFromBase10(value string) *big.Int {
	n, _ := new(big.Int).SetString(value, decimalBase)
	return n
}

type fq [fqLen]uint64

func (e *fq) String() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x", e[5], e[4], e[3], e[2], e[1], e[0])
}

func newFq(str string) (out *fq) {
	words := bigFromBase10(str).Bits()
	numWords := len(words)
	if strconv.IntSize == 64 {
		for i := 0; i < numWords; i++ {
			out[i] = uint64(words[i])
		}
	} else {
		for i := 0; i < numWords; i++ {
			out[i/2] = uint64(words[i]) << uint(32*(i%2))
		}
	}
	return
}
