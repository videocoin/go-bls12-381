package bls12

import "fmt"

// fq2 is an element of Fq² = Fq[X]/(X² − β), where β
// is a quadratic non-residue in Fq with a value of -1.
// See http://eprint.iacr.org/2006/471.pdf - Quadratic extensions.
type fq2 struct{ c0, c1 fq }

func newFq2(c0, c1 fq) fq2 {
	return fq2{
		c0: c0,
		c1: c1,
	}
}

func (fq2 *fq2) String() string {
	return fmt.Sprintf("c0: %s, c1: %s", fq2.c0.String(), fq2.c1.String())
}

func fq2Add(z, x, y *fq2) {
	fqAdd(&z.c0, &x.c0, &y.c0)
	fqAdd(&z.c1, &x.c1, &y.c1)
}

func fq2Sub(z, x, y *fq2) {
	fqSub(&z.c0, &x.c0, &y.c0)
	fqSub(&z.c1, &x.c1, &y.c1)
}

func fq2Mul(z, x, y *fq2) {
	v0, v1 := new(fq), new(fq)
	// v0 = a0b0
	fqMul(v0, &x.c0, &y.c0)
	// v1 = a1b1
	fqMul(v1, &x.c1, &y.c1)
	// c0 = v0 + βv1
	// c0 = v0 - v1
	fqSub(&z.c0, v0, v1)

	// c1 = (a0 + a1)(b0 + b1) − v0 − v1
	// c1 = a0b1 + a1b0
	tmp := new(fq)
	fqMul(&z.c1, &x.c0, &y.c1)
	fqMul(tmp, &x.c1, &y.c0)
	fqAdd(&z.c1, &z.c1, tmp)
}

func fq2Sqr(c, a *fq2) {
	fq2Mul(c, a, a)
}

func fq2Dbl(c, a *fq2) {
	fq2Add(c, a, a)
}
