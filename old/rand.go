package bls12

import (
	"crypto/rand"
	"io"
	"math/big"
)

// randFieldElementIndex returns a random index given the prime that defines the field, q.
func randFieldElementIndex(reader io.Reader) (n *big.Int, err error) {
	return randInt(reader, q)
}

// randInt returns a random scalar between 0 and max.
func randInt(reader io.Reader, max *big.Int) (n *big.Int, err error) {
	for {
		n, err = rand.Int(reader, max)
		if n.Sign() > 0 || err != nil {
			return
		}
	}
}
