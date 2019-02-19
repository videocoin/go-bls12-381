package bls12

import (
	"math/big"
)

var g2Gen = &twistPoint{
	x: fq2{
		c0: g2X0,
		c1: g2X1,
	},
	y: fq2{
		c0: g2Y0,
		c1: g2Y1,
	},
}

// g2 is an abstract cyclic group of the bilinear map
type g2 struct{}

func (g2 *g2) elementAt(n *big.Int) *twistPoint {
	return new(twistPoint).mul(g2Gen, n)
}
