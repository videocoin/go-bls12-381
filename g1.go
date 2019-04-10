package bls12

import (
	"math/big"

	"golang.org/x/crypto/blake2b"
)

var (
	g1Gen = &g1Point{newCurvePoint(g1X, g1Y)}

	g10 = []byte("G1_0")
	g11 = []byte("G1_1")
)

type g1Point struct {
	p *curvePoint
}

func newG1Point() *g1Point {
	return &g1Point{p: new(curvePoint)}
}

func (z *g1Point) ScalarBaseMult(scalar *big.Int) *g1Point {
	return z.ScalarMult(g1Gen, scalar)
}

func (z *g1Point) ScalarMult(x *g1Point, scalar *big.Int) *g1Point {
	z.p.ScalarMult(x.p, scalar)
	return z
}

func (z *g1Point) Add(x, y *g1Point) *g1Point {
	z.p.Add(x.p, y.p)

	return z
}

func (z *g1Point) SetBytes(buf []byte) (*g1Point, error) {
	// See https://github.com/Chia-Network/bls-signatures/blob/master/SPEC.md#hashg1
	h := blake2b.Sum256(buf)
	sum := blake2b.Sum512(append(h[:], g10...))
	bigG0 := new(big.Int).Mod(new(big.Int).SetBytes(sum[:]), q)
	sum = blake2b.Sum512(append(h[:], g11...))
	bigG1 := new(big.Int).Mod(new(big.Int).SetBytes(sum[:]), q)
	fqG0, err := fqMontgomeryFromBig(bigG0)
	if err != nil {
		return nil, err
	}
	fqG1, err := fqMontgomeryFromBig(bigG1)
	if err != nil {
		return nil, err
	}

	return z.Add(g1PointFromFq(&fqG0), g1PointFromFq(&fqG1)).ScalarBaseMult(g1Cofactor), nil
}

func (z *g1Point) Marshal() []byte {
	// TODO
	return []byte("")
}

// g1PointFromFq implements the Shallue and van de Woestijne encoding.
// The point is not guaranteed to be in a particular subgroup.
func g1PointFromFq(t *fq) *g1Point {
	// See https://www.di.ens.fr/~fouque/pub/latincrypt12.pdf - Algorithm 1
	// w = (t^2 + 4u + 1)^(-1) * sqrt(-3) * t
	w, inv := new(fq), new(fq)
	fqMul(w, fqSqrtNeg3, t)
	fqMul(inv, t, t)
	fqAdd(inv, inv, &curveB)
	fqAdd(inv, inv, &fqMont1)
	fqInv(inv, inv)
	fqMul(w, w, inv)

	x, y := new(fq), new(fq)
	for i := 0; i < 3; i++ {
		switch i {
		// x = (sqrt(-3) - 1) / 2 - (w * t)
		case 0:
			fqMul(x, t, w)
			fqSub(x, fqHalfSqrNeg3Minus1, x)
		// x = -1 - x
		case 1:
			fqSub(x, fqNeg1, x)
		// x = 1/w^2 + 1
		case 2:
			fqSqr(x, w)
			fqInv(x, x)
			fqAdd(x, x, &fq1)
		}

		// y^2 = x^3 + 4u
		fqCube(y, x)
		fqAdd(y, y, &curveB)

		// y = sqrt(y2)
		if fqSqrt(y, y) {
			return &g1Point{newCurvePoint(*x, *y)}
		}
	}

	return nil
}
