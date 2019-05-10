package bls12

import (
	"math/big"
)

var g1Gen = &g1Point{curvePoint{*g1X, *g1Y, *fqOne}}

type g1Point struct {
	p curvePoint
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (z *g1Point) ScalarBaseMult(scalar *big.Int) *g1Point {
	return z.ScalarMult(g1Gen, scalar)
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
func (z *g1Point) ScalarMult(x *g1Point, scalar *big.Int) *g1Point {
	z.p.ScalarMult(&x.p, scalar)
	return z
}

// Add returns the sum of (x1,y1) and (x2,y2)
func (z *g1Point) Add(x, y *g1Point) *g1Point {
	z.p.Add(&x.p, &y.p)
	return z
}

// SetBytes uses the Shallue and van de Woestijne encoding.
// The point is guaranteed to be in the subgroup.
func (z *g1Point) SetBytes(buf []byte) *g1Point {
	z.p.SetBytes(buf, g10, g11)
	return z.ScalarMult(z, g1Cofactor)
}

func (z *g1Point) Marshal() []byte {
	return z.p.Marshal()
}

func (z *g1Point) Unmarshal(data []byte) error {
	return z.p.Unmarshal(data)
}
