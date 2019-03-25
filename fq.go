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
	// fqLen is the expected length of a field element
	FqLen            = 6
	fqCompressedLen  = 48
	fqUnompressedLen = 96

	decimalBase = 10
)

var ErrOutOfBounds = errors.New("value is not an element of the finite field of order q")

var (
	// field elements
	fq0 = Fq{0}
	fq1 = Fq{1}

	// field elements in the Montgomery form
	fqMont1, _ = FqMontgomeryFromBase10("1")

	// swEncode helpers
	fqNeg1              = new(Fq)
	fqSqrtNeg3          = &Fq{}
	fqHalfSqrNeg3Minus1 = &Fq{}

	fqQMinus3Over4 []uint64
)

func init() {
	// TODO replace with exact value
	fqNeg(fqNeg1, &fq1)
}

type (
	// fq is an element of the finite field of order q.
	Fq [FqLen]uint64
	// fqLarge is used for storing the basic multiplication result.
	FqLarge [FqLen * 2]uint64
)

// Hex returns the field element in the hexadecimal base
func (fq *Fq) Hex() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x", fq[5], fq[4], fq[3], fq[2], fq[1], fq[0])
}

// String satisfies the Stringer interface.
func (fq *Fq) String() string {
	return fq.Hex()
}

// Bytes returns the absolute value of fq as a big-endian byte slice.
func (fq *Fq) Bytes() (ret []byte) {
	ret = make([]byte, fqCompressedLen)
	for i, fqi := range fq {
		binary.LittleEndian.PutUint64(ret[i+i*8:], fqi)
	}

	return
}

func (fl *FqLarge) Hex() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x", fl[11], fl[10], fl[9], fl[8], fl[7], fl[6], fl[5], fl[4], fl[3], fl[2], fl[1], fl[0])
}

func (fl *FqLarge) String() string {
	return fl.Hex()
}

// IsFieldElement checks if value is within the field bounds.
func IsFieldElement(value *big.Int) bool {
	return (value.Sign() >= 0) && (value.Cmp(q) < 0)
}

func BigFromBase10(str string) *big.Int {
	n, _ := new(big.Int).SetString(str, decimalBase)
	return n
}

// FqFromBase10 converts a base10 value to a field element.
func FqFromBase10(str string) (Fq, error) {
	return FqFromBig(BigFromBase10(str))
}

// FqMontgomeryFromBase10 converts a base10 value to a field element in the Montgomery form.
func FqMontgomeryFromBase10(str string) (Fq, error) {
	return FqMontgomeryFromBig(BigFromBase10(str))
}

// FqFromBig converts a big integer to a field element.
func FqFromBig(value *big.Int) (Fq, error) {
	if !IsFieldElement(value) {
		return Fq{}, ErrOutOfBounds
	}

	fq := Fq{0}
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

	return fq, nil
}

// FqMontgomeryFromBig converts a big integer to a field element in the Montgomery form.
func FqMontgomeryFromBig(value *big.Int) (Fq, error) {
	fieldElement, err := FqFromBig(value)
	if err != nil {
		return Fq{}, err
	}
	montgomeryEncode(&fieldElement, &fieldElement)

	return fieldElement, nil
}

// fqFromHash converts a hash value to a field element.
// See https://golang.org/src/crypto/ecdsa/ecdsa.go?s=1572:1621#L118
func fqFromHash(hash []byte) Fq {
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}
	bigInt := new(big.Int).SetBytes(hash)
	excess := orderBytes*8 - orderBits
	if excess > 0 {
		bigInt.Rsh(bigInt, uint(excess))
	}

	ret, _ := FqFromBig(bigInt)

	return ret
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

// RandFieldElement returns a random element of the field underlying the given
// curve.
func RandFieldElement(reader io.Reader) (n *big.Int, err error) {
	return randInt(reader, q)
}

// montEncode converts the input to Montgomery form.
func montgomeryEncode(c, a *Fq) {
	// See http://home.deib.polimi.it/pelosi/lib/exe/fetch.php?media=teaching:montgomery.pdf page 12/17
	fqMul(c, a, &r2)
}

// montDecode converts the input in the Montgomery form back to
// the standard form.
func montgomeryDecode(c, a *Fq) {
	// See http://home.deib.polimi.it/pelosi/lib/exe/fetch.php?media=teaching:montgomery.pdf page 12/17
	fqMul(c, a, &fq1)
}

// coordinatesFromFq implements the Shallue and van de Woestijne encoding.
// The point is not guaranteed to be in a particular subgroup.
// See https://www.di.ens.fr/~fouque/pub/latincrypt12.pdf
func coordinatesFromFq(t Fq) (x, y Fq) {
	// w = (t^2 + 4u + 1)^(-1) * sqrt(-3) * t
	w, inv := new(Fq), new(Fq)
	fqMul(w, fqSqrtNeg3, &t)
	fqMul(inv, &t, &t)
	fqAdd(inv, inv, &curveB)
	fqAdd(inv, inv, &fqMont1)
	fqInv(inv, inv)
	fqMul(w, w, inv)

	for i := 0; i < 3; i++ {
		switch i {
		// x = (sqrt(-3) - 1) / 2 - (w * t)
		case 0:
			fqMul(&x, &t, w)
			fqSub(&x, fqHalfSqrNeg3Minus1, &x)
		// x = -1 - x
		case 1:
			fqSub(&x, fqNeg1, &x)
		// x = 1/w^2 + 1
		case 2:
			fqSqr(&x, w)
			fqInv(&x, &x)
			fqAdd(&x, &x, &fq1)
		}

		// y^2 = x^3 + 4u
		fqCube(&y, &x)
		fqAdd(&y, &y, &curveB)

		// y = sqrt(y2)
		if fqSqrt(&y, &y) {
			return
		}
	}

	return
}
