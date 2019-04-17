package bls12

var fq2One = &fq2{fqMont1, fq0}

type (
	// fq2 is an element of Fq² = Fq[X]/(X² − β), where β
	// is a quadratic non-residue in Fq with a value of -1.
	// See http://eprint.iacr.org/2006/471.pdf - Quadratic extensions.
	fq2      struct{ c0, c1 fq }
	fq2Large [2]fqLarge
)

func newFq2(c0, c1 fq) fq2 {
	return fq2{
		c0: c0,
		c1: c1,
	}
}

func (fq *fq2) IsOne() bool {
	return fq.c0 == fq2One.c0 && fq.c1 == fq2One.c1
}

func (z *fq2) SetOne() *fq2 {
	z.c0.SetOne()
	z.c1.SetZero()
	return z
}

func (z *fq2) SetZero() *fq2 {
	z.c0.SetZero()
	z.c1.SetZero()
	return z
}

func (z *fq2) Add(x, y *fq2) *fq2 {
	fqAdd(&z.c0, &x.c0, &y.c0)
	fqAdd(&z.c1, &x.c1, &y.c1)
	return z
}

func (z *fq2) Sub(x, y *fq2) *fq2 {
	fqSub(&z.c0, &x.c0, &y.c0)
	fqSub(&z.c1, &x.c1, &y.c1)
	return z
}

// Karatsuba method
// note: there's room for optimization (multiplications + reductions).
func (z *fq2) Mul(x, y *fq2) *fq2 {
	mult := new(fq2)
	// v0 = a0b0
	// v1 = a1b1
	// c0 = v0 + βv1
	// c0 = v0 - v1
	v0, v1 := new(fq), new(fq)
	fqMul(v0, &x.c0, &y.c0)
	fqMul(v1, &x.c1, &y.c1)
	fqSub(&mult.c0, v0, v1)

	// c1 = (a0 + a1)(b0 + b1) − v0 − v1
	// c1 = (a0 + a1)(b0 + b1) − (v0 + v1)
	// note: c1 = a0b1 + a1b0 = expensive 2M
	a, b, v := new(fq), new(fq), new(fq)
	fqAdd(a, &x.c0, &x.c1)
	fqAdd(b, &y.c0, &y.c1)
	fqMul(&mult.c1, a, b)
	fqAdd(v, v0, v1)
	fqSub(&mult.c1, &mult.c1, v)

	z.c0, z.c1 = mult.c0, mult.c1

	return z
}

// TODO Karatsuba
func (z *fq2) Sqr(x *fq2) *fq2 {
	return z.Mul(x, x)
}

func (z *fq2) Dbl(x *fq2) *fq2 {
	return z.Add(x, x)
}

func (x *fq2) Equal(y *fq2) bool {
	// TODO a.c0.Equal form
	return (x.c0 == y.c0) && (x.c1 == y.c1)
}

func (z *fq2) Inv(x *fq2) *fq2 {
	// t0 = a.c0^2
	// t1 = a.c1^2
	// t0 = t0 + t1
	// t0 = 1/t0
	// c.c0 = a.c0 * t0
	// c.c1 = - a.c1 * t0
	t0, t1 := new(fq), new(fq)
	fqSqr(t0, &x.c0)
	fqSqr(t1, &x.c1)
	fqAdd(t0, t0, t1)
	fqInv(t0, t0)
	fqMul(&z.c0, &x.c0, t0)
	fqMul(&z.c1, &x.c1, t0)
	fqNeg(&z.c1, &z.c1)

	return z
}
