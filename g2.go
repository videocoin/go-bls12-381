package bls12

import "math/big"

var (
	G2          = new(g2)
	g2Generator = newTwistPoint(newFq2(g2X0, g2X1), newFq2(g2Y0, g2Y1))
)

type g2Point = twistPoint

type g2 struct{}

func (g2 *g2) Element(index *big.Int) *twistPoint {
	return new(g2Point).ScalarMult(g2Generator, index)
}
