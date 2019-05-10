package bls12

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"strconv"
)

const (
	fqLen       = 6
	fqByteLen   = 48
	decimalBase = 10
)

var (
	fqZero        = fq{}
	fqOne, _      = new(fq).SetString("1")
	fqOneStandard = fq{1}
)

// fq is an element of the finite field of order q.
// fq operates, internally, on the montgomery form but it's possible to
// represent the element on the standard form by using the decoding methods
// available or using the struct literal. Note that the user is responsible for
// making sure that the montgomery form is used whenever required.
type fq [fqLen]uint64

// IsOne reports whether x is equal to 1.
func (x *fq) IsOne() bool {
	return *x == *fqOne
}

// String implements the Stringer interface.
func (x *fq) String() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x", x[5], x[4], x[3], x[2], x[1], x[0])
}

// Set sets z to x and returns z.
func (z *fq) Set(x *fq) *fq {
	*z = *x
	return z
}

// SetOne sets z to 0 and returns z.
func (z *fq) SetZero() *fq {
	*z = fqZero
	return z
}

// SetOne sets z to 1 and returns z.
func (z *fq) SetOne() *fq {
	return z.Set(fqOne)
}

// MontgomeryEncode converts z to the Montgomery form and returns z.
// See http://home.deib.polimi.it/pelosi/lib/exe/fetch.php?media=teaching:montgomery.pdf page 12/17
func (z *fq) MontgomeryEncode(x *fq) *fq {
	fqMul(z, x, r2)
	return z
}

// MontgomeryDecode converts z back to the standard form and returns z.
// See http://home.deib.polimi.it/pelosi/lib/exe/fetch.php?media=teaching:montgomery.pdf page 12/17
func (z *fq) MontgomeryDecode(x *fq) *fq {
	fqMul(z, x, &fqOneStandard)
	return z
}

// SetString sets z to the Montgomery value of s, interpreted in the decimal
// base, and returns z and a boolean indicating success.
func (z *fq) SetString(s string) (*fq, bool) {
	k, valid := bigFromBase10(s)
	if !valid {
		return nil, false
	}

	return z.SetInt(k)
}

// SetInt sets z to the Mongomery value of x and returns z and a boolean
// indicating success. The integer must be within field bounds for success. If
// the operation failed, the value of z is undefined but the returned value is
// nil.
func (z *fq) SetInt(x *big.Int) (*fq, bool) {
	if !isFieldElement(x) {
		return nil, false
	}

	fq := fq{0}
	words := x.Bits()
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

	return z.Set(&fq).MontgomeryEncode(z), true
}

// SetUint64 sets z to the value of x and returns z.
func (z *fq) SetUint64(x uint64) *fq {
	return z.Set(&fq{x}).MontgomeryEncode(z)
}

// Bytes returns the absolute value of fq as a big-endian byte slice.
func (x *fq) Bytes() []byte {
	ret := make([]byte, fqByteLen)
	for i, xi := range x {
		binary.BigEndian.PutUint64(ret[fqByteLen-(i+1)*8:], xi)
	}
	return ret
}

// Int returns the corresponding big integer.
func (x *fq) Int() *big.Int {
	var words []big.Word

	if strconv.IntSize == 64 {
		words = make([]big.Word, 0, fqLen)
		for _, word := range x {
			words = append(words, big.Word(word))
		}
	} else {
		numWords := fqLen * 2
		words = make([]big.Word, 0, numWords)
		for i := 0; i < numWords; i++ {
			words = append(words, big.Word(uint32((x[i/2])>>uint(32*(i%2)))))
		}
	}

	return new(big.Int).SetBits(words)
}

// fqLarge is used during the multiplication.
type fqLarge [fqLen * 2]uint64

// String implements the Stringer interface.
func (x *fqLarge) String() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x", x[11], x[10], x[9], x[8], x[7], x[6], x[5], x[4], x[3], x[2], x[1], x[0])
}

// isFieldElement reports whether the value is within field bounds.
func isFieldElement(value *big.Int) bool {
	return (value.Sign() >= 0) && (value.Cmp(q) < 0)
}

// randInt returns a random scalar between 0 and max.
func randInt(reader io.Reader, max *big.Int) (n *big.Int, err error) {
	for {
		n, err = rand.Int(reader, max)
		if n.Sign() > 0 || err != nil {
			return
		}
	}
}

// randFieldElement returns a random scalar between 0 and q.
func randFieldElement(reader io.Reader) (*big.Int, error) {
	return randInt(reader, q)
}

func bigFromBase10(s string) (*big.Int, bool) {
	return new(big.Int).SetString(s, decimalBase)
}
