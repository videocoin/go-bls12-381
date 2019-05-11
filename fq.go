package bls12

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strconv"
)

const (
	fqLen       = 6
	fqByteLen   = 48
	decimalBase = 10
	wordSize    = 64
)

var errOutOfBounds = errors.New("fq: value must be within the bounds of the field")

var (
	fqZero        = fq{}
	fqOne, _      = new(fq).SetString("1")
	fqMinusOne, _ = new(fq).SetString("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559786")
	// fqOneStarndard is the value by which to multiply field elements to map
	// them back to the standard form.
	fqOneStandard        = fq{1}
	minusThreeOverFour64 = []uint64{0xee7fbfffffffeaaa, 0x7aaffffac54ffff, 0xd9cc34a83dac3d89, 0xd91dd2e13ce144af, 0x92c6e9ed90d2eb35, 0x680447a8e5ff9a6}
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

// Equal reports whether x is equal to y.
func (x *fq) Equal(y *fq) bool {
	return *x == *y
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

func bigFromBase10(s string) (*big.Int, bool) {
	return new(big.Int).SetString(s, decimalBase)
}

// SetString sets z to the Montgomery value of s, interpreted in the decimal
// base, and returns z and a boolean indicating success.
func (z *fq) SetString(s string) (*fq, error) {
	k, valid := bigFromBase10(s)
	if !valid {
		return nil, nil
	}

	return z.SetInt(k)
}

// SetInt sets z to the Mongomery value of x and returns z and a boolean
// indicating success. The integer must be within field bounds for success. If
// the operation failed, the value of z is undefined but the returned value is
// nil.
func (z *fq) SetInt(x *big.Int) (*fq, error) {
	if !isFieldElement(x) {
		return nil, errOutOfBounds
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

	return z.Set(&fq).MontgomeryEncode(z), nil
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

func fqInv(z, x *fq) {
	fqExp(z, x, qMinusTwo[:])
}

// See https://www.coursera.org/lecture/mathematical-foundations-cryptography/square-and-multiply-ty62K
func fqExp(z *fq, x *fq, y []uint64) {
	b := *x
	ret := new(fq).SetOne()
	for _, word := range y {
		for j := uint(0); j < wordSize; j++ {
			if (word & (1 << j)) != 0 {
				fqMul(ret, ret, &b)
			}
			fqMul(&b, &b, &b)
		}
	}

	z.Set(ret)
}

// See https://eprint.iacr.org/2012/685.pdf - Algorithm 2; q ≡ 3 (mod 4)
func fqSqrt(x, a *fq) bool {
	a1, a0 := new(fq), new(fq)
	fqExp(a1, a, minusThreeOverFour64)

	fqMul(a0, a1, a1)
	fqMul(a0, a0, a)

	if a0.Equal(fqMinusOne) {
		return false
	}

	fqMul(x, a1, a)

	return true
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
