package bls12

import (
	"math/big"
)

var g2Gen = &g2Point{*newTwistPoint(fq2{*g2X0, *g2X1}, fq2{*g2Y0, *g2Y1})}

type g2Point struct {
	p twistPoint
}

// Set sets z to the value of x and returns z.
func (z *g2Point) Set(x *g2Point) *g2Point {
	if z != x {
		z.p.Set(&x.p)
	}
	return z
}

// Equal reports whether x is equal to y.
func (x *g2Point) Equal(y *g2Point) bool {
	return x.p == y.p
}

// Add sets z to the sum x+y and returns z.
func (z *g2Point) Add(x, y *g2Point) *g2Point {
	z.p.Add(&x.p, &y.p)
	return z
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (z *g2Point) ScalarBaseMult(scalar *big.Int) *g2Point {
	return z.ScalarMult(g2Gen, scalar)
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
func (z *g2Point) ScalarMult(x *g2Point, scalar *big.Int) *g2Point {
	z.p.ScalarMult(&x.p, scalar)
	return z
}

func (z *g2Point) ToAffine() *g2Point {
	z.p.ToAffine()
	return z
}
