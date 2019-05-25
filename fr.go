package bls12

import "math/big"

const (
	frLen = 4
)

// fr is an element of the finite field of order r.
// fr operates, internally, on the montgomery form.
type fr [frLen]uint64

func (z *fr) Set(x *fr) *fr {
	// TODO
	return nil
}

func (z *fr) SetInt(x *big.Int) (*fr, error) {
	// TODO
	return nil, nil
}

func (z *fr) SetUint64(x uint64) *fr {
	//TODO
	return nil
}

// Int returns the corresponding big integer.
func (x *fr) Int() *big.Int {
	/*
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
	*/
	return nil
}

// See https://www.coursera.org/lecture/mathematical-foundations-cryptography/square-and-multiply-ty62K
func frExp(z *fr, x *fr, y []uint64) {
	/*
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
	*/
}

func frInv(z, x *fr) {
	frExp(z, x, rMinusTwo[:])
}

/*

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
*/
