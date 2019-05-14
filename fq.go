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
	fqMinusOne, _ = new(fq).SetString("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559786")
)

// fq is an element of the finite field of order q.
// fq operates, internally, on the montgomery form but it's possible to
// represent the element on the standard form by using the decoding methods
// available or using the struct literal. Note that the user is responsible for
// making sure that the montgomery form is used whenever required.
type fq [fqLen]uint64

// IsOne reports whether x is equal to 1.
func (x *fq) IsOne() bool {
	return *x == *new(fq).SetUint64(1)
}

// Set sets z to x and returns z.
func (z *fq) Set(x *fq) *fq {
	*z = *x
	return z
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
	fqMul(z, x, &fq{1})
	return z
}

// String implements the Stringer interface.
func (x *fq) String() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x", x[5], x[4], x[3], x[2], x[1], x[0])
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

	return z.MontgomeryEncode(&fq), nil
}

// SetUint64 sets z to the value of x and returns z.
func (z *fq) SetUint64(x uint64) *fq {
	return z.MontgomeryEncode(&fq{x})
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
	xDecoded := new(fq).MontgomeryDecode(x)

	if strconv.IntSize == 64 {
		words = make([]big.Word, 0, fqLen)
		for _, word := range xDecoded {
			words = append(words, big.Word(word))
		}
	} else {
		numWords := fqLen * 2
		words = make([]big.Word, 0, numWords)
		for i := 0; i < numWords; i++ {
			words = append(words, big.Word(uint32((xDecoded[i/2])>>uint(32*(i%2)))))
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
	ret := new(fq).SetUint64(1)
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

// fqSqrt sets z to the square root of x, and returns a boolean.
// If it exists, x satisfying x 2 = a, false otherwise.
// See https://eprint.iacr.org/2012/685.pdf - Algorithm 2; q ≡ 3 (mod 4)
// TODO replace x, a with z, x
// TODO desc bool
func fqSqrt(z, x *fq) bool {
	x0, x1 := new(fq), new(fq)
	fqExp(x1, x, []uint64{0xee7fbfffffffeaaa, 0x7aaffffac54ffff, 0xd9cc34a83dac3d89, 0xd91dd2e13ce144af, 0x92c6e9ed90d2eb35, 0x680447a8e5ff9a6})
	fqMul(x0, x1, x1)
	fqMul(x0, x0, x)
	if (*x0 == fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}) {
		return false
	}
	fqMul(z, x1, x)

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
