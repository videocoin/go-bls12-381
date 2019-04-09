package bls12

import "math/big"

var g1Gen = &g1Point{newCurvePoint(g1X, g1Y)}

type g1Point struct {
	p *curvePoint
}

func newG1Point() *g1Point {
	return &g1Point{p: new(curvePoint)}
}

func (z *g1Point) ScalarBaseMult(scalar *big.Int) *g1Point {
	return z.ScalarMult(g1Gen, scalar)
}

func (z *g1Point) ScalarMult(x *g1Point, scalar *big.Int) *g1Point {
	z.p.ScalarMult(x.p, scalar)
	return z
}

func (z *g1Point) SetBytes(buf []byte) *g1Point {
	// TODO
	return &g1Point{}
}

func (z *g1Point) Marshal() []byte {
	// TODO
	return []byte("")
}
