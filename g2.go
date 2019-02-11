package bls12

import (
	"crypto/rand"
	"io"
	"math/big"
)

var (
	g2Generator = &twistPoint{}
)

func randN(reader io.Reader) (n *big.Int, err error) {
	for {
		n, err = rand.Int(reader, Q)
		if n.Sign() > 0 || err != nil {
			return
		}
	}
}

// g2 is an abstract cyclic group. The zero value is suitable for use as the
// output of an operation, but cannot be used as an input.
type g2 struct{}

// randomG2 returns n and g2^n where n is a random, non-zero number read from the reader.
func randomG2(reader io.Reader) (*big.Int, *twistPoint, error) {
	n, err := randN(reader)
	if err != nil {
		return nil, nil, err
	}

	return n, new(g2).scalarBaseMult(n), nil
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (g *g2) scalarBaseMult(scalar *big.Int) *twistPoint {
	return new(twistPoint).mul(g2Generator, scalar)
}
