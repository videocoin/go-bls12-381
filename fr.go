package bls12

import (
	"math/big"
	"strconv"
)

const frLen = 4

// fr is an element of the finite field of order r.
// fr operates, internally, on the montgomery form.
type fr [frLen]uint64

func (z *fr) Set(x *fr) *fr {
	*z = *x
	return z
}

func (z *fr) SetInt(x *big.Int) (*fr, error) {
	if !isFieldElement(x, r) {
		return nil, errOutOfBounds
	}

	fr := fr{0}
	words := x.Bits()
	numWords := len(words)
	if strconv.IntSize == 64 {
		for i := 0; i < numWords; i++ {
			fr[i] = uint64(words[i])
		}
	} else {
		for i := 0; i < numWords; i++ {
			fr[i/2] |= uint64(words[i]) << uint(32*(i%2))
		}
	}

	return z.MontgomeryEncode(&fr), nil
}

func (z *fr) SetUint64(x uint64) *fr {
	return z.MontgomeryEncode(&fr{x})
}

// Int returns the corresponding big integer.
func (x *fr) Int() *big.Int {
	var words []big.Word
	xDecoded := new(fr).MontgomeryDecode(x)

	if strconv.IntSize == 64 {
		words = make([]big.Word, 0, frLen)
		for _, word := range xDecoded {
			words = append(words, big.Word(word))
		}
	} else {
		numWords := frLen * 2
		words = make([]big.Word, 0, numWords)
		for i := 0; i < numWords; i++ {
			words = append(words, big.Word(uint32((xDecoded[i/2])>>uint(32*(i%2)))))
		}
	}

	return new(big.Int).SetBits(words)
}

// See https://www.coursera.org/lecture/mathematical-foundations-cryptography/square-and-multiply-ty62K
func frExp(z *fr, x *fr, y []uint64) {
	b := *x
	ret := new(fr).SetUint64(1)
	for _, word := range y {
		for j := uint(0); j < wordSize; j++ {
			if (word & (1 << j)) != 0 {
				frMul(ret, ret, &b)
			}
			frMul(&b, &b, &b)
		}
	}

	z.Set(ret)
}

func frInv(z, x *fr) {
	// TODO confirm rMinusTwo
	frExp(z, x, rMinusTwo[:])
}

// MontgomeryEncode converts z to the Montgomery form and returns z.
// See http://home.deib.polimi.it/pelosi/lib/exe/fetch.php?media=teaching:montgomery.pdf page 12/17
func (z *fr) MontgomeryEncode(x *fr) *fr {
	frMul(z, x, rR2)
	return z
}

// MontgomeryDecode converts z back to the standard form and returns z.
// See http://home.deib.polimi.it/pelosi/lib/exe/fetch.php?media=teaching:montgomery.pdf page 12/17
func (z *fr) MontgomeryDecode(x *fr) *fr {
	frMul(z, x, &fr{1})
	return z
}
