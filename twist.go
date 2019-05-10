package bls12

import (
	"math/big"
)

// twistPoint is a curve point in the elliptic curve's twist
// over an extension field Fq². T = z1²
// To obtain the full speed of pairings on Weierstrass curves it is useful
//to represent a point by (X1 : Y1 : Z1 : T1) with T1 = Z²
type twistPoint struct {
	x, y, z, t fq2
}

func newTwistPoint(x, y fq2) *twistPoint {
	return &twistPoint{
		x: x,
		y: y,
		z: *fq2One,
		t: *fq2One,
	}
}

func (c *twistPoint) Set(a *twistPoint) *twistPoint {
	c.x, c.y, c.z, c.t = a.x, a.y, a.z, a.t
	return c
}

func (c *twistPoint) Add(a, b *twistPoint) *twistPoint {
	if a.IsInfinity() {
		return c.Set(b)
	}
	if b.IsInfinity() {
		return c.Set(a)
	}

	if a.Equal(b) {
		return c.Double(a) // faster than Add
	}

	// See https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	z1z1 := new(fq2).Sqr(&a.z)
	z2z2 := new(fq2).Sqr(&b.z)

	u1 := new(fq2).Mul(&a.x, z2z2)
	u2 := new(fq2).Mul(&b.x, z1z1)

	s1 := new(fq2).Mul(&a.y, &b.z)
	s1.Mul(s1, z2z2)
	s2 := new(fq2).Mul(&b.y, &a.z)
	s2.Mul(s2, z1z1)

	h := new(fq2).Sub(u2, u1)
	i := new(fq2).Add(h, h)
	i.Sqr(i)
	j := new(fq2).Mul(h, i)
	r := new(fq2).Sub(s2, s1)
	r.Add(r, r)
	v := new(fq2).Mul(u1, i)

	p, t0, t1 := new(twistPoint), new(fq2), new(fq2)
	t0.Add(v, v)
	t0.Add(t0, j)
	p.x.Sqr(r)
	p.x.Sub(&p.x, t0)

	t0.Add(s1, s1)
	t0.Mul(t0, j)
	t1.Sub(v, &p.x)
	t1.Mul(t1, r)
	p.y.Sub(t1, t0)

	p.z.Add(&a.z, &b.z)
	p.z.Sqr(&p.z)
	t0.Add(z1z1, z2z2)
	p.z.Sub(&p.z, t0)
	p.z.Mul(&p.z, h)

	return c.Set(p)
}

func (tp *twistPoint) Double(p *twistPoint) *twistPoint {
	// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
	a, b, c, d, e, f := new(fq2), new(fq2), new(fq2), new(fq2), new(fq2), new(fq2)
	a.Sqr(&p.x)
	b.Sqr(&p.y)
	c.Sqr(b)
	d.Mul(&p.x, b)
	d.Add(d, d)
	d.Add(d, d)
	e.Add(a, a)
	e.Add(e, a)
	f.Sqr(e)

	x, y, z, t0 := new(fq2), new(fq2), new(fq2), new(fq2)
	x.Add(d, d)
	x.Sub(f, x)
	t0.Add(c, c)
	t0.Add(t0, t0)
	t0.Add(t0, t0)
	y.Sub(d, x)
	y.Mul(y, e)
	y.Sub(y, t0)
	z.Mul(&p.y, &p.z)
	z.Add(z, z)

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

// See https://www.sciencedirect.com/topics/computer-science/affine-coordinate - Jacobian Projective Points
func (tp *twistPoint) ToAffine() *twistPoint {
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
