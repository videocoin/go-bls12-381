package bls12

import "math/big"

var (
	g2Generator = newTwistPoint(newFq2(g2X0, g2X1), newFq2(g2Y0, g2Y1))

	G2 = new(g2)
)

type g2 struct{}

func (g2 *g2) element(index *big.Int) *twistPoint {
	return new(twistPoint).mul(g2Generator, index)
}
