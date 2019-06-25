package bls12

import (
	"math/big"
)

var g1Gen = &G1Point{curvePoint{
	x: *g1X,
	y: *g1Y,
	z: *new(fq).SetUint64(1),
}}

type G1Point struct {
	p curvePoint
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (z *G1Point) ScalarBaseMult(scalar *big.Int) *G1Point {
	return z.ScalarMult(g1Gen, scalar)
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
func (z *G1Point) ScalarMult(x *G1Point, scalar *big.Int) *G1Point {
	z.p.ScalarMult(&x.p, scalar)
	return z
}

// Add returns the sum of (x1,y1) and (x2,y2)
func (z *G1Point) Add(x, y *G1Point) *G1Point {
	z.p.Add(&x.p, &y.p)
	return z
}

// HashToPoint uses the Shallue and van de Woestijne encoding.
// The point is guaranteed to be in the subgroup.
func (z *G1Point) HashToPoint(buf []byte) *G1Point {
	z.p.HashToPoint(buf, g10, g11)
	return z.ScalarMult(z, g1Cofactor)
}

func (z *G1Point) ToAffine() *G1Point {
	z.p.ToAffine()
	return z
}

func (z *G1Point) Marshal() []byte {
	return z.p.Marshal()
}

func (z *G1Point) Unmarshal(data []byte) error {
	return z.p.Unmarshal(data)
}
