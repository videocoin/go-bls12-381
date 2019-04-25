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
	fqLen     = 6
	fqByteLen = 48

	decimalBase = 10
)

var errOutOfBounds = errors.New("value is not an element of the finite field of order q")

var (
	fq0        = fq{0}
	fq1        = fq{1}
	fqMont1, _ = fqMontgomeryFromBase10("1")
)

type (
	// fq is an element of the finite field of order q.
	fq [fqLen]uint64
	// fqLarge is used for storing the basic multiplication result.
	fqLarge [fqLen * 2]uint64
)

// Hex returns the field element in the hexadecimal base
func (fq *fq) Hex() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x", fq[5], fq[4], fq[3], fq[2], fq[1], fq[0])
}

func (z *fq) Set(x *fq) *fq {
	*z = *x
	return z
}

func (fq *fq) SetOne() *fq {
	*fq = fqMont1
	return fq
}

func (fq *fq) SetZero() *fq {
	*fq = fq0
	return fq
}

// String satisfies the Stringer interface.
func (fq *fq) String() string {
	return fq.Hex()
}

// Bytes returns the absolute value of fq as a big-endian byte slice.
func (fq *fq) Bytes() []byte {
	ret := make([]byte, fqByteLen)

	for i, fqi := range fq {
		binary.BigEndian.PutUint64(ret[fqByteLen-(i+1)*8:], fqi)
	}

	return ret
}

func (fl *fqLarge) Hex() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x%16.16x", fl[11], fl[10], fl[9], fl[8], fl[7], fl[6], fl[5], fl[4], fl[3], fl[2], fl[1], fl[0])
}

func (fl *fqLarge) String() string {
	return fl.Hex()
}

// IsFieldElement checks if value is within the field bounds.
func IsFieldElement(value *big.Int) bool {
	return (value.Sign() >= 0) && (value.Cmp(q) < 0)
}

func bigFromBase10(str string) *big.Int {
	n, _ := new(big.Int).SetString(str, decimalBase)
	return n
}

// fqFromBase10 converts a base10 value to a field element.
func fqFromBase10(str string) (fq, error) {
	return fqFromBig(bigFromBase10(str))
}

// fqMontgomeryFromBase10 converts a base10 value to a field element in the Montgomery form.
func fqMontgomeryFromBase10(str string) (fq, error) {
	return fqMontgomeryFromBig(bigFromBase10(str))
}

// fqFromBig converts a big integer to a field element.
func fqFromBig(value *big.Int) (fq, error) {
	if !IsFieldElement(value) {
		return fq{}, errOutOfBounds
	}

	return fqFromFqBig(value), nil
}

// fqFromFqBig converts the big integer representing the
// field element to the field element representation.
func fqFromFqBig(value *big.Int) fq {
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

	return fq
}

// fqMontgomeryFromBig converts a big integer to a field element in the Montgomery form.
func fqMontgomeryFromBig(value *big.Int) (fq, error) {
	fieldElement, err := fqFromBig(value)
	if err != nil {
		return fq{}, err
	}
	montgomeryEncode(&fieldElement, &fieldElement)

	return fieldElement, nil
}

// fqFromFqMontgomery decodes a field element in the montgomery form.
func fqFromFqMontgomery(fq fq) fq {
	montgomeryDecode(&fq, &fq)
	return fq
}

// fqFromHash converts a hash value to a field element.
// See https://golang.org/src/crypto/ecdsa/ecdsa.go?s=1572:1621#L118
func fqFromHash(hash []byte) fq {
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}
	bigInt := new(big.Int).SetBytes(hash)
	excess := orderBytes*8 - orderBits
	if excess > 0 {
		bigInt.Rsh(bigInt, uint(excess))
	}

	ret, _ := fqFromBig(bigInt)

	return ret
}

// FqToBig converts a field element to a big integer.
func fqToBig(fq fq) *big.Int {
	var words []big.Word

	if strconv.IntSize == 64 {
		words = make([]big.Word, 0, fqLen)
		for _, word := range fq {
			words = append(words, big.Word(word))
		}
	} else {
		numWords := fqLen * 2
		words = make([]big.Word, 0, numWords)
		for i := 0; i < numWords; i++ {
			words = append(words, big.Word(uint32((fq[i/2])>>uint(32*(i%2)))))
		}
	}

	return new(big.Int).SetBits(words)
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
/*
func RandFieldElement(reader io.Reader) (Fq, error) {
	elem, err := randInt(reader, q)
	if err != nil {
		return Fq{}, err
	}

	// TODO verification is not necessary out of bounds
	return FqMontgomeryFromBig(elem)
}
*/

func randFieldElement(reader io.Reader) (*big.Int, error) {
	return randInt(reader, q)
}

// montEncode converts the input to Montgomery form.
func montgomeryEncode(c, a *fq) {
	// See http://home.deib.polimi.it/pelosi/lib/exe/fetch.php?media=teaching:montgomery.pdf page 12/17
	fqMul(c, a, &r2)
}

// montDecode converts the input in the Montgomery form back to
// the standard form.
func montgomeryDecode(c, a *fq) {
	// See http://home.deib.polimi.it/pelosi/lib/exe/fetch.php?media=teaching:montgomery.pdf page 12/17
	fqMul(c, a, &fq1)
}

func (fq *fq) IsOne() bool {
	return *fq == fqMont1
}
