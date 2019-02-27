package bls12

import (
	"math/big"
)

var twistB = twistPoint{}

// twist point is an eliptic curve(y²=x³+4(u+1)) point over the finite field Fq².
type twistPoint struct {
	x, y, z fq2
}

func (tp *twistPoint) isInfinity() bool {
	return tp.z.isZero()
}

// Add sets tp to the sum a+b and returns c.

func (tp *twistPoint) mul(p *twistPoint, scalar *big.Int) *twistPoint {
	// See https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add
	for i := scalar.BitLen(); i > 0; i-- {
		tp.double(tp)
		if scalar.Bit(i) != 0 {
			tp.add(tp, p)
		}
	}

	return tp
}
