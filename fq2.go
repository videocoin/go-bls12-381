package bls12

import "fmt"

// fq2 is an element of Fq² = Fq[X]/(X² − β), where β
// is a quadratic non-residue in Fq with a value of -1.
// See http://eprint.iacr.org/2006/471.pdf - Quadratic extensions.
type (
	fq2 struct{ c0, c1 Fq }

	fq2Large [2]FqLarge
)

func newFq2(c0, c1 Fq) fq2 {
	return fq2{
		c0: c0,
		c1: c1,
	}
}

func (fq2 *fq2) String() string {
	return fmt.Sprintf("c0: %s, c1: %s", fq2.c0.String(), fq2.c1.String())
}

func fq2Add(z, x, y *fq2) {
	FqAdd(&z.c0, &x.c0, &y.c0)
	FqAdd(&z.c1, &x.c1, &y.c1)
}

func fq2Sub(z, x, y *fq2) {
	FqSub(&z.c0, &x.c0, &y.c0)
	FqSub(&z.c1, &x.c1, &y.c1)
}

func fq2BasicMul(z *fq2Large, x *fq2, y *fq2) {
	// Karatsuba method
	v0, v1 := new(FqLarge), new(FqLarge)
	// v0 = a0b0
	FqBasicMul(v0, &x.c0, &y.c0)
	// v1 = a1b1
	FqBasicMul(v1, &x.c1, &y.c1)

	// c0 = v0 + βv1
	// c0 = v0 - v1
	fqLargeSub(&z[0], v0, v1)

	// c1 = (a0 + a1)(b0 + b1) − v0 − v1
	// c1 = (a0 + a1)(b0 + b1) − (v0 + v1)
	// c1 = a0b1 + a1b0
	t0, t1 := new(Fq), new(Fq)
	sum := new(FqLarge)
	FqAdd(t0, &x.c0, &x.c1)
	FqAdd(t1, &y.c0, &y.c1)
	FqBasicMul(&z[1], t0, t1)
	fqLargeAdd(sum, v0, v1)
	fqLargeSub(&z[1], &z[1], sum)
}

func fq2REDC(z *fq2, x *fq2Large) {
	FqREDC(&z.c0, &x[0])
	FqREDC(&z.c1, &x[1])
}

func fq2Mul(z, x, y *fq2) {
	large := new(fq2Large)
	fq2BasicMul(large, x, y)
	fq2REDC(z, large)
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
