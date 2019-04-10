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
		// cheaper
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

	x, y, z, t0, t1 := new(fq2), new(fq2), new(fq2), new(fq2), new(fq2)
	fq2Dbl(t0, v)
	fq2Add(t0, t0, j)
	fq2Sqr(x, r)
	fq2Sub(x, x, t0)

	fq2Dbl(t0, s1)
	fq2Mul(t0, t0, j)
	fq2Sub(t1, v, x)
	fq2Mul(t1, t1, r)
	fq2Sub(y, t1, t0)

	fq2Add(z, &a.z, &b.z)
	fq2Sqr(z, z)
	fq2Add(t0, z1z1, z2z2)
	fq2Sub(z, z, t0)
	fq2Mul(z, z, h)

	tp.x, tp.y, tp.z = *x, *y, *z

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

	x, y, z, t0 := new(fq2), new(fq2), new(fq2), new(fq2)
	fq2Dbl(x, d)
	fq2Sub(x, f, x)
	fq2Dbl(t0, c)
	fq2Dbl(t0, t0)
	fq2Dbl(t0, t0)
	fq2Sub(y, d, x)
	fq2Mul(y, y, e)
	fq2Sub(y, y, t0)
	fq2Mul(z, &p.y, &p.z)
	fq2Dbl(z, z)

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
	return fq2Equal(&tp.x, &p.x) && fq2Equal(&tp.y, &p.y) && fq2Equal(&tp.z, &p.z)
}

/*
func (tp *twistPoint) String() string {
	tp.MakeAffine()
	x, y := fq2Decode(&tp.x), fq2Decode(&tp.y)

	return fmt.Sprintf("x: %s, y: %s", x.String(), y.String())
}
*/

func (tp *twistPoint) ToAffine() *twistPoint {
	// See https://www.sciencedirect.com/topics/computer-science/affine-coordinate - Jacobian Projective Points
	if tp.z.IsOne() {
		return tp
	}

	zInv, zInvSqr, zInvCube := new(fq2), new(fq2), new(fq2)
	fq2Inv(zInv, &tp.z)
	fq2Mul(zInvSqr, zInv, zInv)
	fq2Mul(zInvCube, zInvSqr, zInv)
	fq2Mul(&tp.x, &tp.x, zInvSqr)
	fq2Mul(&tp.y, &tp.y, zInvCube)
	tp.z = fq2{fqMont1, fq0}

	return tp
}
