package bls12

import (
	"fmt"
	"math/big"
)

var g2Gen = &g2Point{newTwistPoint(newFq2(g2X0, g2X1), newFq2(g2Y0, g2Y1))}

type g2Point struct {
	p *twistPoint
}

func newG2Point() *g2Point {
	return &g2Point{p: new(twistPoint)}
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (z *g2Point) ScalarBaseMult(scalar *big.Int) *g2Point {
	return z.ScalarMult(g2Gen, scalar)
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
func (z *g2Point) ScalarMult(x *g2Point, scalar *big.Int) *g2Point {
	z.p.ScalarMult(x.p, scalar)
	return z
}

func (z *g2Point) Equal(x *g2Point) bool {
	return z.p.Equal(x.p)
}

// Add returns the sum of (x1,y1) and (x2,y2).
func (z *g2Point) Add(x, y *g2Point) *g2Point {
	z.p.Add(x.p, y.p)
	return z
}

func (z *g2Point) Set(x *g2Point) *g2Point {
	if z != x {
		z.p.Set(x.p)
	}
	return z
}

func (z *g2Point) ToAffine() *g2Point {
	z.p.ToAffine()
	return z
}

func (z *g2Point) String() string {
	return fmt.Sprintf("%v", z.p)
}
