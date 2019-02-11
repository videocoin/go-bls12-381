package bls12

import (
	"math/big"
)

var (
	twistB = twistPoint{}
)

// twist point implements the eliptic curve y2 = x3 + 4(u + 1) over GF(fq2)
type twistPoint struct {
	x, y, z fq2
}

// Add sets tp to the sum a+b and returns c.
func (tp *twistPoint) Add(a, b *twistPoint) *twistPoint {
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
	return tp
}

func (tp *twistPoint) Double(a *twistPoint) *twistPoint {
	if a == nil {
		*a = *tp
	}
	// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l

	return tp
}

func (tp *twistPoint) mul(point *twistPoint, scalar *big.Int) *twistPoint {
	return tp
}
