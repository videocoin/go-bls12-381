package bls12

import (
	"math/big"
)

var (
	g1Gen = &g1Point{newCurvePoint(g1X, g1Y)}
	g10   = []byte("G1_0")
	g11   = []byte("G1_1")
)

type g1Point struct {
	p *curvePoint
}

func newG1Point() *g1Point {
	return &g1Point{p: new(curvePoint)}
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (z *g1Point) ScalarBaseMult(scalar *big.Int) *g1Point {
	return z.ScalarMult(g1Gen, scalar)
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
func (z *g1Point) ScalarMult(x *g1Point, scalar *big.Int) *g1Point {
	z.p.ScalarMult(x.p, scalar)
	return z
}

// Add returns the sum of (x1,y1) and (x2,y2)
func (z *g1Point) Add(x, y *g1Point) *g1Point {
	z.p.Add(x.p, y.p)
	return z
}

func (z *g1Point) SetBytes(buf []byte) *g1Point {
	z.p.SetBytes(buf)
	return z.ScalarMult(z, g1Cofactor) // map to g1
}

// TODO
func (z *g1Point) Marshal() []byte {
	return []byte("")
}
