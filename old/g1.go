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
	g11 = []byte("G1_1")
)

// g1 is an abstract cyclic group of the bilinear map
type g1 struct{}

func (g1 *g1) elementAt(n *big.Int) *curvePoint {
	return new(curvePoint).mul(g1Gen, n)
}

func hashToG1(hash []byte) *curvePoint {
	hasher, _ := blake2b.New512(nil)

	//t0 := hasher.Sum(nil)
	hasher.Write(hash)
	hasher.Write(g10)
	var t0 fq

	//t1 := hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(hash)
	hasher.Write(g11)
	var t1 fq

	swEncodeToG1(t0)
	swEncodeToG1(t1)

	return &curvePoint{}
}

// swEncodeG1 Shallue–van de Woestijne encoding to BN.
func swEncodeToG1(t fq) *curvePoint {
	// See https://www.di.ens.fr/~fouque/pub/latincrypt12.pdf -  "Algorithm 1"
	if t.isZero() {
		// point at infinity
	}

	w := new(fq)
	fqSqr(w, &t)
	fqAdd(w, w, &curveB)
	fqAdd(w, w, &fq1)
	fqInv(w, w)
	fqMul(w, &t, w)
	fqMul(w, fqSqrtNeg3, w)

	x, y := new(fq), new(fq)
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			// x1 ← (−1 + √−3)/2 − t
			fqMul(x, &t, w)
			fqSub(x, fqHalfSqrNeg3Minus1, x)
		case 1:
			// x2 ← −1 − x
			fqSub(x, fqNeg1, x)
		case 2:
			// x3 ← 1 + 1/w
			fqSqr(x, w)
			fqInv(x, x)
			fqAdd(x, x, &fq1)
		}

		// y
		fqSqr(y, x)
		fqMul(y, y, y)
		fqAdd(y, y, &curveB)
		fqSqrt(y, y)

		if y != nil {
			return newCurvePoint(*x, *y)
		}
	}

	return nil
}
