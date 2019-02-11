package bls12

import "math/big"

var (
	curveB = newFq("4")
)

// curvePoint implements the elliptic curve y²=x³+3 over GF(fq)
type curvePoint struct {
	x, y, z fq
}

// Add sets cp to the sum a+b and returns c.
func (cp *curvePoint) Add(a, b *curvePoint) *curvePoint {
	/*
		if a == nil {
			*a = *c
		}

		if a.IsInfinity() {
			c.Set(b)
			return c
		}
		if b.IsInfinity() {
			c.Set(a)
			return c
		}

		// See https://hyperelliptic.org/EFD/g1p/auto-code/shortw/jacobian-3/addition/mmadd-2007-bl.op3

	*/
	return cp
}

func (cp *curvePoint) Double(a *curvePoint) *curvePoint {
	if a == nil {
		*a = *cp
	}
	// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l

	return cp
}

func (cp *curvePoint) scalarBaseMult(scalar *big.Int) *curvePoint {
	return &curvePoint{}
}
