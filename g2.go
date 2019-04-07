package bls12

import (
	"fmt"
	"math/big"
)

var (
	G2    = new(g2)
	g2Gen = &g2Point{newTwistPoint(newFq2(g2X0, g2X1), newFq2(g2Y0, g2Y1))}
)

type g2Point struct {
	p *twistPoint
}

func NewG2Point() *g2Point {
	return &g2Point{
		p: new(twistPoint),
	}
}

func (z *g2Point) ScalarMult(x *g2Point, scalar *big.Int) *g2Point {
	z.p.ScalarMult(x.p, scalar)
	return z
}

func (z *g2Point) Equal(x *g2Point) bool {
	return z.p.Equal(x.p)
}

func (z *g2Point) Add(x, y *g2Point) *g2Point {
	z.p.Add(x.p, y.p)
	return z
}

func (z *g2Point) Set(x *g2Point) *g2Point {
	z.p.Set(x.p)
	return z
}

func (z *g2Point) String() string {
	return fmt.Sprintf("%v", z.p)
}

type g2 struct{}

func (g2 *g2) ScalarBaseMult(index *big.Int) *g2Point {
	return NewG2Point().ScalarMult(g2Gen, index)
}
