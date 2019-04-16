package bls12

import (
	"fmt"
	"math/big"
)

var g2Gen = &g2Point{newTwistPoint(newFq2(g2X0, g2X1), newFq2(g2Y0, g2Y1))}

type g2Point struct {
	p *twistPoint
}

func newG2Point() *g2Point {
	return &g2Point{p: new(twistPoint)}
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (z *g2Point) ScalarBaseMult(scalar *big.Int) *g2Point {
	return z.ScalarMult(g2Gen, scalar)
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
func (z *g2Point) ScalarMult(x *g2Point, scalar *big.Int) *g2Point {
	z.p.ScalarMult(x.p, scalar)
	return z
}

func (z *g2Point) Equal(x *g2Point) bool {
	return z.p.Equal(x.p)
}

// Add returns the sum of (x1,y1) and (x2,y2).
func (z *g2Point) Add(x, y *g2Point) *g2Point {
	z.p.Add(x.p, y.p)
	return z
}

func (z *g2Point) Set(x *g2Point) *g2Point {
	if z != x {
		z.p.Set(x.p)
	}
	return z
}

func (z *g2Point) ToAffine() *g2Point {
	z.p.ToAffine()
	return z
}

func (z *g2Point) String() string {
	return fmt.Sprintf("%v", z.p)
}

// DoubleLine sets z to the double of x and f to the line function result and returns the pair (z, f).
// See https://arxiv.org/pdf/0904.0854v3.pdf - Doubling on curves with a4 = 0
// todo: q must be affine?
func (z *g2Point) DoubleLine(r *g2Point, q *g1Point, f *fq12) (*g2Point, *fq12) {
	// R ← [2]R
	// note: there's a faster way to compute the doubling (2m + 5s instead of 1m +
	// 7s) but line functions make use of T1 = Z².
	// todo: benchmark variable allocation
	ret := new(twistPoint)
	a := new(fq2).Sqr(&r.p.x)
	b := new(fq2).Sqr(&r.p.y)
	c := new(fq2).Sqr(b)
	d := new(fq2).Dbl(new(fq2).Sub(new(fq2).Sqr(new(fq2).Add(&r.p.x, b)), new(fq2).Add(a, c)))
	e := new(fq2).Add(new(fq2).Dbl(a), a)
	g := new(fq2).Sqr(e)
	ret.x.Sub(g, new(fq2).Dbl(d))
	ret.y.Sub(new(fq2).Mul(e, new(fq2).Sub(d, &ret.x)), new(fq2).Dbl(new(fq2).Dbl(new(fq2).Dbl(c))))
	ret.z.Sub(new(fq2).Sqr(new(fq2).Add(&r.p.y, &r.p.z)), new(fq2).Add(b, &r.p.t))
	ret.t.Sqr(&ret.z)

	return nil, &fq12{}
}

func (z *g2Point) AddLine(r *g2Point, q *g1Point) (*g2Point, *fq12) {
	return nil, &fq12{}
}
