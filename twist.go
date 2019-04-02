package bls12

import (
	"math/big"
)

// twistPoint is a curve point in the elliptic curve's twist
// over an extension field Fq².
type twistPoint struct {
	x, y, z fq2
}

func newTwistPoint(x, y fq2) *twistPoint {
	return &twistPoint{
		x: x,
		y: y,
		z: fq2{FqMont1, Fq0},
	}
}

// note: room for optimization if numbers are equal (cheaper to use double instead of add)
func (tp *twistPoint) Add(a, b *twistPoint) *twistPoint {
	if a.Equal(b) {
		return tp.Double(a)
	}

	// See https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	z1z1, z2z2 := new(fq2), new(fq2)
	fq2Sqr(z1z1, &a.z)
	fq2Sqr(z2z2, &b.z)

	u1, u2 := new(fq2), new(fq2)
	fq2Mul(u1, &a.x, z2z2)
	fq2Mul(u2, &b.x, z1z1)

	s1, s2 := new(fq2), new(fq2)
	fq2Mul(s1, &a.y, &b.z)
	fq2Mul(s1, s1, z2z2)
	fq2Mul(s2, &b.y, &a.z)
	fq2Mul(s2, s2, z1z1)

	h, i, j, r, v := new(fq2), new(fq2), new(fq2), new(fq2), new(fq2)
	fq2Sub(h, u2, u1)
	fq2Dbl(i, h)
	fq2Sqr(i, i)
	fq2Mul(j, h, i)
	fq2Sub(r, s2, s1)
	fq2Dbl(r, r)
	fq2Mul(v, u1, i)

	t0, t1 := new(fq2), new(fq2)
	fq2Dbl(t0, v)
	fq2Sqr(&tp.x, r)
	fq2Sub(&tp.x, &tp.x, j)
	fq2Sub(&tp.x, &tp.x, t0)
	fq2Dbl(t0, s1)
	fq2Mul(t0, t0, j)
	fq2Mul(t1, r, &tp.x)
	fq2Mul(&tp.y, r, v)
	fq2Sub(&tp.y, &tp.y, t1)
	fq2Sub(&tp.y, &tp.y, t0)
	fq2Mul(&tp.z, &a.z, &b.z)
	fq2Sqr(&tp.z, &tp.z)
	fq2Sub(&tp.z, &tp.z, z1z1)
	fq2Sub(&tp.z, &tp.z, z2z2)
	fq2Mul(&tp.z, &tp.z, h)

	return tp
}

func (tp *twistPoint) Double(p *twistPoint) *twistPoint {
	// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
	// D=4*X1*B
	a, b, c, d, e, f := new(fq2), new(fq2), new(fq2), new(fq2), new(fq2), new(fq2)
	fq2Sqr(a, &p.x)
	fq2Sqr(b, &p.y)
	fq2Sqr(c, b)
	fq2Mul(d, &p.x, b)
	fq2Dbl(d, d)
	fq2Dbl(d, d)
	fq2Dbl(e, a)
	fq2Add(e, e, a)
	fq2Sqr(f, e)

	fq2Dbl(&tp.x, d)
	fq2Sub(&tp.x, f, &tp.x)

	t0 := new(fq2)
	fq2Dbl(t0, c)
	fq2Dbl(t0, t0)
	fq2Dbl(t0, t0)
	fq2Sub(&tp.y, d, &tp.x)
	// fixme: fq2Mul(&tp.y, e, &tp.y) // Error (references)
	// fixme: fq2Mul(&tp.y, &tp.y, e) // Error (references)
	t1 := tp.y
	fq2Mul(&tp.y, e, &t1)
	fq2Sub(&tp.y, &tp.y, t0)

	fq2Mul(&tp.z, &p.y, &p.z)
	fq2Dbl(&tp.z, &tp.z)

	return tp
}

func (tp *twistPoint) ScalarMult(p *twistPoint, scalar *big.Int) *twistPoint {
	// See https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add
	q := new(twistPoint)
	for i := scalar.BitLen() - 1; i >= 0; i-- {
		q.Double(q)
		if scalar.Bit(i) == 1 {
			q.Add(q, p)
		}
	}
	*tp = *q

	return tp
}

func (tp *twistPoint) Equal(p *twistPoint) bool {
	return fq2Equal(&tp.x, &p.x) && fq2Equal(&tp.y, &p.y) && fq2Equal(&tp.z, &p.z)
}

// twistSubGroup is a cyclic group of the elliptic curve twist.
type twistSubGroup struct {
	generator *twistPoint
}

func newTwistSubGroup(gen *twistPoint) *twistSubGroup {
	return &twistSubGroup{
		generator: gen,
	}
}
