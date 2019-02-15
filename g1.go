package bls12

import (
	"math/big"

	"golang.org/x/crypto/blake2b"
)

var (
	g1Gen = &curvePoint{
		x: g1X,
		y: g1Y,
	}

	g10 = []byte("G1_0")
	g11 = []byte("G1_0")
)

type g1 struct{}

// g1 is an abstract cyclic group
func (g1 *g1) elementAt(n *big.Int) *curvePoint {
	return new(curvePoint).mul(g1Gen, n)
}

func hashToG1(hash []byte) *curvePoint {
	hasher, _ := blake2b.New512(nil)

	hasher.Write(hash)
	hasher.Write(g10)
	//t0 := hasher.Sum(nil)

	hasher.Reset()
	hasher.Write(hash)
	hasher.Write(g11)
	//t1 := hasher.Sum(nil)

	return &curvePoint{}
}

//  swEncG1 Shallue–van de Woestijne encoding to BN curves
func swEncG1(t fq) *curvePoint {
	// See https://www.di.ens.fr/~fouque/pub/latincrypt12.pdf
	return &curvePoint{}
}
