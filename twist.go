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
		z: *fq2One,
	}
}

func (tp *twistPoint) Set(p *twistPoint) *twistPoint {
	tp.x, tp.y, tp.z = p.x, p.y, p.z

	return tp
}

// note: room for optimization if numbers are equal (cheaper to use double instead of add)
func (tp *twistPoint) Add(a, b *twistPoint) *twistPoint {
	// TODO is infinity confirm
	if a.IsInfinity() {
		return tp.Set(b)
	}

	// TODO is infinity confirm
	if b.IsInfinity() {
		return tp.Set(a)
	}

	if a.Equal(b) {
		return tp.Double(a) // faster than Add
	}

	// See https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	z1z1, z2z2 := new(fq2), new(fq2)
	z1z1.Sqr(&a.z)
	z2z2.Sqr(&b.z)

	u1, u2 := new(fq2), new(fq2)
	u1.Mul(&a.x, z2z2)
	u2.Mul(&b.x, z1z1)

	s1, s2 := new(fq2), new(fq2)
	s1.Mul(&a.y, &b.z)
	s1.Mul(s1, z2z2)
	s2.Mul(&b.y, &a.z)
	s2.Mul(s2, z1z1)

	h, i, j, r, v := new(fq2), new(fq2), new(fq2), new(fq2), new(fq2)
	h.Sub(u2, u1)
	i.Dbl(h)
	i.Sqr(i)
	j.Mul(h, i)
	r.Sub(s2, s1)
	r.Dbl(r)
	v.Mul(u1, i)

	x, y, z, t0, t1 := new(fq2), new(fq2), new(fq2), new(fq2), new(fq2)
	t0.Dbl(v)
	t0.Add(t0, j)
	x.Sqr(r)
	x.Sub(x, t0)

	t0.Dbl(s1)
	t0.Mul(t0, j)
	t1.Sub(v, x)
	t1.Mul(t1, r)
	y.Sub(t1, t0)

	z.Add(&a.z, &b.z)
	z.Sqr(z)
	t0.Add(z1z1, z2z2)
	z.Sub(z, t0)
	z.Mul(z, h)

	tp.x, tp.y, tp.z = *x, *y, *z

	return tp
}

func (tp *twistPoint) Double(p *twistPoint) *twistPoint {
	// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
	a, b, c, d, e, f := new(fq2), new(fq2), new(fq2), new(fq2), new(fq2), new(fq2)
	a.Sqr(&p.x)
	b.Sqr(&p.y)
	c.Sqr(b)
	d.Mul(&p.x, b)
	d.Dbl(d)
	d.Dbl(d)
	e.Dbl(a)
	e.Add(e, a)
	f.Sqr(e)

	x, y, z, t0 := new(fq2), new(fq2), new(fq2), new(fq2)
	x.Dbl(d)
	x.Sub(f, x)
	t0.Dbl(c)
	t0.Dbl(t0)
	t0.Dbl(t0)
	y.Sub(d, x)
	y.Mul(y, e)
	y.Sub(y, t0)
	z.Mul(&p.y, &p.z)
	z.Dbl(z)

	tp.x, tp.y, tp.z = *x, *y, *z

	return tp
}

// TODO confirm logic behind it
func (c *twistPoint) IsInfinity() bool {
	return c.z == fq2{}
}

func (tp *twistPoint) ScalarMult(p *twistPoint, scalar *big.Int) *twistPoint {
	// See https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add
	q := &twistPoint{}
	// fixme: BitLen or BitLen - 1
	for i := scalar.BitLen() - 1; i >= 0; i-- {
		q.Double(q)
		if scalar.Bit(i) == 1 {
			q.Add(q, p)
		}
	}

	return tp.Set(q)
}

func (tp *twistPoint) Equal(p *twistPoint) bool {
	return tp.x.Equal(&p.x) && tp.y.Equal(&p.y) && tp.z.Equal(&p.z)
}

func (tp *twistPoint) ToAffine() *twistPoint {
	// See https://www.sciencedirect.com/topics/computer-science/affine-coordinate - Jacobian Projective Points
	if tp.z.IsOne() {
		return tp
	}

	zInv, zInvSqr, zInvCube := new(fq2), new(fq2), new(fq2)
	zInv.Inv(&tp.z)
	zInvSqr.Sqr(zInv)
	zInvCube.Mul(zInvSqr, zInv)
	tp.x.Mul(&tp.x, zInvSqr)
	tp.y.Mul(&tp.y, zInvCube)
	tp.z.SetOne()

	return tp
}
