package bls12

import (
	"math/big"
)

var g2Gen = &G2Point{*newTwistPoint(fq2{*g2X0, *g2X1}, fq2{*g2Y0, *g2Y1})}

type G2Point struct {
	p twistPoint
}

// Set sets z to the value of x and returns z.
func (z *G2Point) Set(x *G2Point) *G2Point {
	if z != x {
		z.p.Set(&x.p)
	}
	return z
}

// Equal reports whether x is equal to y.
func (x *G2Point) Equal(y *G2Point) bool {
	return x.p == y.p
}

// Add sets z to the sum x+y and returns z.
func (z *G2Point) Add(x, y *G2Point) *G2Point {
	z.p.Add(&x.p, &y.p)
	return z
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (z *G2Point) ScalarBaseMult(scalar *big.Int) *G2Point {
	return z.ScalarMult(g2Gen, scalar)
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
func (z *G2Point) ScalarMult(x *G2Point, scalar *big.Int) *G2Point {
	z.p.ScalarMult(&x.p, scalar)
	return z
}

func (z *G2Point) ToAffine() *G2Point {
	z.p.ToAffine()
	return z
}

func (z *G2Point) HashToPoint(buf []byte) *G2Point {
	// TODO review
	return z.HashToPointWithDomain(buf, 0)
}

// HashToPointWithDomain uses the Shallue and van de Woestijne encoding.
// The point is guaranteed to be in the subgroup.
func (z *G2Point) HashToPointWithDomain(buf []byte, domain uint64) *G2Point {
	//z.p.HashToPoint(buf, g10, g11)
	return z.ScalarMult(z, g2Cofactor)
}
