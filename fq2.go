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

/*
func (fq2 *fq2) String() string {
	return fmt.Sprintf("%s + %s*X\n", fq2.c0.String(), fq2.c1.String())
}
*/

func fq2Add(z, x, y *fq2) {
	fqAdd(&z.c0, &x.c0, &y.c0)
	fqAdd(&z.c1, &x.c1, &y.c1)
}

func fq2Sub(z, x, y *fq2) {
	fqSub(&z.c0, &x.c0, &y.c0)
	fqSub(&z.c1, &x.c1, &y.c1)
}

// note: there's room for optimization (multiplications + reductions).
func fq2Mul(z, x, y *fq2) {
	// Karatsuba method
	// v0 = a0b0
	// v1 = a1b1
	// c0 = v0 + βv1
	// c0 = v0 - v1
	v0, v1 := new(fq), new(fq)
	fqMul(v0, &x.c0, &y.c0)
	fqMul(v1, &x.c1, &y.c1)
	fqSub(&z.c0, v0, v1)

	// c1 = (a0 + a1)(b0 + b1) − v0 − v1
	// c1 = (a0 + a1)(b0 + b1) − (v0 + v1)
	// note: c1 = a0b1 + a1b0 = expensive 2M
	a, b, v := new(fq), new(fq), new(fq)
	fqAdd(a, &x.c0, &x.c1)
	fqAdd(b, &y.c0, &y.c1)
	fqMul(&z.c1, a, b)
	fqAdd(v, v0, v1)
	fqSub(&z.c1, &z.c1, v)
}

func fq2Sqr(c, a *fq2) {
	fq2Mul(c, a, a)
}

func fq2Dbl(c, a *fq2) {
	fq2Add(c, a, a)
}

func fq2Equal(a, b *fq2) bool {
	return (a.c0 == b.c0) && (a.c1 == b.c1)
}

func fq2Decode(a *fq2) *fq2 {
	fq2 := new(fq2)
	montgomeryDecode(&fq2.c0, &a.c0)
	montgomeryDecode(&fq2.c1, &a.c1)

	return fq2
}

func fq2Inv(a, b *fq2) {
	// TODO
}
