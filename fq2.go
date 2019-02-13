package bls12

// fq2 is an element of Fq2 = Fq[X]/(X2 − β), where β
// is a quadratic non-residue in Fq with a value of -1.
// fq2 is represented as as α0 + α1X, where αi ∈ Fq.
type fq2 struct {
	c0, c1 fq
}

func (elm *fq2) isZero() bool {
	return elm.c0.isZero() && elm.c1.isZero()
}

func fq2Add(c, a, b *fq2) {
	fqAdd(&c.c0, &a.c0, &b.c0)
	fqAdd(&c.c1, &a.c1, &b.c1)
}

func fq2Sub(c, a, b *fq2) {
	fqSub(&c.c0, &a.c0, &b.c0)
	fqSub(&c.c1, &a.c1, &b.c1)
}

func fq2Mul(c, a, b *fq2) {
	// See http://eprint.iacr.org/2006/471.pdf - Quadratic extensions
	t0, t1 := new(fq), new(fq)
	fqMul(t0, &a.c0, &b.c0)
	fqMul(t1, &a.c1, &b.c1)
	fqSub(&c.c0, t0, t1)

	fqMul(t0, &a.c0, &b.c1)
	fqMul(t1, &a.c1, &b.c0)
	fqAdd(&c.c1, t0, t1)
}

func fq2Sqr(c, a *fq2) {
	fq2Mul(c, a, a)
}

func fq2Dbl(c, a *fq2) {
	fq2Add(c, a, a)
}
