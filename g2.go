package bls12

import (
	"crypto/rand"
	"io"
	"math/big"
)

var g2Generator = &twistPoint{
	x: fq2{
		c0: g2X0,
		c1: g2X1,
	},
	y: fq2{
		c0: g2Y0,
		c1: g2Y1,
	},
}

func randN(reader io.Reader) (n *big.Int, err error) {
	for {
		n, err = rand.Int(reader, q)
		if n.Sign() > 0 || err != nil {
			return
		}
	}
}

// g2 is an abstract cyclic group.
type g2 struct{}

// randomG2 returns n and g2^n where n is a random, non-zero number read from the reader.
func randomG2(reader io.Reader) (*big.Int, *twistPoint, error) {
	n, err := randN(reader)
	if err != nil {
		return nil, nil, err
	}

	return n, new(g2).scalarBaseMult(n), nil
}

// scalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (g *g2) scalarBaseMult(scalar *big.Int) *twistPoint {
	return new(twistPoint).mul(g2Generator, scalar)
}
