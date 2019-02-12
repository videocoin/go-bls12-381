package bls12

import "math/big"

var (
	curveB = newFq(bigFromBase10("4"))
)

// curvePoint implements the elliptic curve y²=x³+3 over GF(fq)
type curvePoint struct {
	x, y, z fq
}

func (cp *curvePoint) set(p *curvePoint) {
	cp.x = p.x
	cp.y = p.y
	cp.z = p.z
}

// Add sets cp to the sum a+b and returns c.
func (cp *curvePoint) add(a, b *curvePoint) *curvePoint {
	/*
		if a == nil {
			*a = *c
		}

		if a.IsInfinity() {
			c.Set(b)
			return c
		}
		if b.IsInfinity() {
			c.Set(a)
			return c
		}

		// See https://hyperelliptic.org/EFD/g1p/auto-code/shortw/jacobian-3/addition/mmadd-2007-bl.op3

	*/
	return cp
}

func (cp *curvePoint) double(p *curvePoint) *curvePoint {
	// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
	a, b, c, d, e, f, temp := new(fq), new(fq), new(fq), new(fq), new(fq), new(fq), new(fq)

	// a
	fqSqr(a, &p.x)

	// b
	fqSqr(b, &p.y)

	// c
	fqSqr(c, b)

	// d
	fqAdd(d, &p.x, b)
	fqSqr(d, d)
	fqSub(d, d, a)
	fqSub(d, d, c)
	fqDbl(d, d)

	// e
	fqDbl(e, a)
	fqAdd(e, e, a)

	// f
	fqSqr(f, e)

	// x3
	fqDbl(&cp.x, d)
	fqSub(&cp.x, f, &cp.x)

	// y3
	fqAdd(temp, c, c)
	fqAdd(temp, temp, temp)
	fqDbl(temp, temp)
	fqSub(&cp.y, d, &cp.x)
	fqMul(&cp.y, e, &cp.y)
	fqSub(&cp.y, &cp.y, temp)

	// z3
	fqMul(&cp.z, &p.y, &p.z)
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
