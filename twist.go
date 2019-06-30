package bls12

import (
	"math/big"
)

var (
	// Values taken from the execution of https://eprint.iacr.org/2019/403.pdf - A The isogeny maps.
	iso3XNum = [11]fq{}
	iso3XDen = [11]fq{}
	iso3YNum = [11]fq{}
	iso3YDen = [11]fq{}
)

// twistPoint is a curve point in the elliptic curve's twist over an extension
// field Fq². T = z1². To obtain the full speed of pairings on Weierstrass
// curves it is useful to represent a point by (X1 : Y1 : Z1 : T1) with T1 = Z²
type twistPoint struct {
	x, y, z, t fq2
}

func newTwistPoint(x, y fq2) *twistPoint {
	return &twistPoint{
		x: x,
		y: y,
		z: fq2{c0: *new(fq).SetUint64(1)},
		t: fq2{c0: *new(fq).SetUint64(1)},
	}
}

// Set sets c to the value of a and returns c.
func (c *twistPoint) Set(a *twistPoint) *twistPoint {
	c.x, c.y, c.z, c.t = a.x, a.y, a.z, a.t
	return c
}

// Equal reports whether a is equal to b.
func (a *twistPoint) Equal(b *twistPoint) bool {
	return *a == *b
}

// IsInfinity reports whether the point is at infinity.
func (a *twistPoint) IsInfinity() bool {
	return a.z == fq2{}
}

// Add sets c to the sum a+b and returns c.
func (c *twistPoint) Add(a, b *twistPoint) *twistPoint {
	if a.IsInfinity() {
		return c.Set(b)
	}
	if b.IsInfinity() {
		return c.Set(a)
	}

	// faster than Add
	if a.Equal(b) {
		return c.Double(a)
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

// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
func (c *twistPoint) Double(a *twistPoint) *twistPoint {
	d, e, f, g, h, i := new(fq2), new(fq2), new(fq2), new(fq2), new(fq2), new(fq2)
	d.Sqr(&a.x)
	e.Sqr(&a.y)
	f.Sqr(e)
	g.Mul(&a.x, e)
	g.Add(g, g)
	g.Add(g, g)
	h.Add(d, d)
	h.Add(h, d)
	i.Sqr(h)

	p, t0 := new(twistPoint), new(fq2)
	p.x.Add(g, g)
	p.x.Sub(i, &p.x)
	t0.Add(f, f)
	t0.Add(t0, t0)
	t0.Add(t0, t0)
	p.y.Sub(g, &p.x)
	p.y.Mul(&p.y, h)
	p.y.Sub(&p.y, t0)
	p.z.Mul(&a.y, &a.z)
	p.z.Add(&p.z, &p.z)

	return c.Set(p)
}

// ScalarMult returns b*(Ax,Ay) where b is a number in big-endian form.
// See https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add.
func (c *twistPoint) ScalarMult(a *twistPoint, b *big.Int) *twistPoint {
	p := &twistPoint{}
	for i := b.BitLen() - 1; i >= 0; i-- {
		p.Double(p)
		if b.Bit(i) == 1 {
			p.Add(p, a)
		}
	}

	return c.Set(p)
}

// ToAffine sets a to its affine value and returns a.
// See https://www.sciencedirect.com/topics/computer-science/affine-coordinate - Jacobian Projective Points
func (a *twistPoint) ToAffine() *twistPoint {
	if (a.z.c0 == *new(fq).SetUint64(1)) && (a.z.c1 == fq{}) {
		return a
	}

	zInv, zInvSqr, zInvCube := new(fq2), new(fq2), new(fq2)
	zInv.Inv(&a.z)
	zInvSqr.Sqr(zInv)
	zInvCube.Mul(zInvSqr, zInv)
	a.x.Mul(&a.x, zInvSqr)
	a.y.Mul(&a.y, zInvCube)
	a.z.SetOne()
	a.t.SetOne()

	return a
}

// TODO
// SetBytes sets c to the twist point that results from the given slice of bytes
// and returns c. The point is not guaranteed to be in a particular subgroup.
// See https://github.com/ethereum/eth2.0-specs/blob/dev/specs/bls_signature.md.
/*
func (c *curvePoint) SetBytes(buf []byte, ref0 []byte, ref1 []byte) *curvePoint {
	h := blake2b.Sum256(buf)
	sum := blake2b.Sum512(append(h[:], ref0...))
	t0 := new(big.Int)
	g10, _ := new(fq).SetInt(t0.Mod(t0.SetBytes(sum[:]), r))
	sum = blake2b.Sum512(append(h[:], ref1...))
	g11, _ := new(fq).SetInt(t0.Mod(t0.SetBytes(sum[:]), r))

	return c.Add(new(curvePoint).SWEncode(g10), new(curvePoint).SWEncode(g11))
}
*/
