package bls12

import (
	"math/big"
)

var (
	twistB = twistPoint{}
)

// twist point implements the eliptic curve y2 = x3 + 4(u + 1) over GF(fq2)
type twistPoint struct {
	x, y, z fq2
}

func (tp *twistPoint) set(p *twistPoint) {
	tp.x = p.x
	tp.y = p.y
	tp.z = p.z
}

// Add sets tp to the sum a+b and returns c.
func (tp *twistPoint) add(a, b twistPoint) *twistPoint {
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
	return tp
}

func (tp *twistPoint) double(p *twistPoint) *twistPoint {
	// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
	a, b, c, d, e, f, temp := new(fq2), new(fq2), new(fq2), new(fq2), new(fq2), new(fq2), new(fq2)

	// a
	fq2Sqr(a, &p.x)

	// b
	fq2Sqr(b, &p.y)

	// c
	fq2Sqr(c, b)

	// d
	fq2Add(d, &p.x, b)
	fq2Sqr(d, d)
	fq2Sub(d, d, a)
	fq2Sub(d, d, c)
	fq2Dbl(d, d)

	// e
	fq2Dbl(e, a)
	fq2Add(e, e, a)

	// f
	fq2Sqr(f, e)

	// x3
	fq2Dbl(&tp.x, d)
	fq2Sub(&tp.x, f, &tp.x)

	// y3
	fq2Add(temp, c, c)
	fq2Add(temp, temp, temp)
	fq2Dbl(temp, temp)
	fq2Sub(&tp.y, d, &tp.x)
	fq2Mul(&tp.y, e, &tp.y)
	fq2Sub(&tp.y, &tp.y, temp)

	// z3
	fq2Mul(&tp.z, &tp.y, &tp.z)
	fq2Add(&tp.z, &tp.z, &tp.z)

	return tp
}

func (tp *twistPoint) mul(p *twistPoint, scalar *big.Int) *twistPoint {
	// See https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add
	q := new(twistPoint)
	for i := scalar.BitLen(); i > 0; i-- {
		//q.double(q)
		if scalar.Bit(i) != 0 {
			//q.add(q, p)
		}
	}
	tp.set(q)

	return tp
}
