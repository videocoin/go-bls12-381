package bls12

import (
	"math/big"
)

var curveB = newFq(bigFromBase10("4"))

// curvePoint implements the elliptic curve y²=x³+3 over GF(fq)
type curvePoint struct {
	x, y, z fq
}

func (cp *curvePoint) isInfinity() bool {
	return cp.z.isZero()
}

func (cp *curvePoint) set(p *curvePoint) {
	cp.x = p.x
	cp.y = p.y
	cp.z = p.z
}

// Add sets cp to the sum a+b and returns c.
func (cp *curvePoint) add(a, b *curvePoint) *curvePoint {
	// See https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	if a.isInfinity() {
		return a
	}
	if b.isInfinity() {
		return b
	}

	z1z1, z2z2 := new(fq), new(fq)
	fqSqr(z1z1, &a.z)
	fqSqr(z2z2, &b.z)

	u1, u2 := new(fq), new(fq)
	fqMul(u1, &a.x, z2z2)
	fqMul(u2, &b.x, z1z1)

	s1, s2 := new(fq), new(fq)
	fqMul(s1, &a.y, &b.z)
	fqMul(s1, s1, z2z2)
	fqMul(s2, &b.y, &a.z)
	fqMul(s2, s2, z1z1)

	h, i, j, r, v := new(fq), new(fq), new(fq), new(fq), new(fq)
	fqSub(h, u2, u1)
	fqDbl(i, h)
	fqSqr(i, i)
	fqMul(j, h, i)
	fqSub(r, s2, s1)
	fqDbl(r, r)
	fqMul(v, u1, i)

	t0, t1 := new(fq), new(fq)
	fqDbl(t0, v)
	fqSqr(&cp.x, r)
	fqSub(&cp.x, &cp.x, j)
	fqSub(&cp.x, &cp.x, t0)
	fqDbl(t0, s1)
	fqMul(t0, t0, j)
	fqMul(t1, r, &cp.x)
	fqMul(&cp.y, r, v)
	fqSub(&cp.y, &cp.y, t1)
	fqSub(&cp.y, &cp.y, t0)
	fqMul(&cp.z, &a.z, &b.z)
	fqSqr(&cp.z, &cp.z)
	fqSub(&cp.z, &cp.z, z1z1)
	fqSub(&cp.z, &cp.z, z2z2)
	fqMul(&cp.z, &cp.z, h)

	return cp
}

func (cp *curvePoint) double(p *curvePoint) *curvePoint {
	// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
	a, b, c, d, e, f := new(fq), new(fq), new(fq), new(fq), new(fq), new(fq)

	fqSqr(a, &p.x)
	fqSqr(b, &p.y)
	fqSqr(c, b)
	fqAdd(d, &p.x, b)
	fqSqr(d, d)
	fqSub(d, d, a)
	fqSub(d, d, c)
	fqDbl(d, d)
	fqDbl(e, a)
	fqAdd(e, e, a)
	fqSqr(f, e)

	t0 := new(fq)
	fqDbl(&cp.x, d)
	fqSub(&cp.x, f, &cp.x)
	fqAdd(t0, c, c)
	fqAdd(t0, t0, t0)
	fqDbl(t0, t0)
	fqSub(&cp.y, d, &cp.x)
	fqMul(&cp.y, e, &cp.y)
	fqSub(&cp.y, &cp.y, t0)
	fqMul(&cp.z, &cp.y, &cp.z)
	fqAdd(&cp.z, &cp.z, &cp.z)

	return cp
}

func (cp *curvePoint) mul(p *curvePoint, scalar *big.Int) *curvePoint {
	// See https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add
	q := new(curvePoint)
	for i := scalar.BitLen(); i > 0; i-- {
		//q.double(q)
		if scalar.Bit(i) != 0 {
			//q.add(q, p)
		}
	}
	cp.set(q)

	return cp
}
