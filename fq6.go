package bls12

// fq6 is an element of Fq6 = Fq2[Y]/(Y3 − γ), where γ
// is a quadratic non-residue in Fq and γ = √β is a
// cubic non-residue in Fq2 with a value of X + 1.
// fq6 is represented as c0 + c1Y + c2Y2, where ci ∈ Fq2.
// See https://eprint.iacr.org/2006/471.pdf - "6.2 Cubic over quadratic"
type fq6 struct {
	c0, c1, c2 fq2
}

func fq6Add(c, a, b *fq6) {
	fq2Add(&c.c0, &a.c0, &b.c0)
	fq2Add(&c.c1, &a.c1, &b.c1)
	fq2Add(&c.c2, &a.c2, &b.c2)
}

func fq6Sub(c, a, b *fq6) {
	fq2Sub(&c.c0, &a.c0, &b.c0)
	fq2Sub(&c.c1, &a.c1, &b.c1)
	fq2Sub(&c.c2, &a.c2, &b.c2)
}

func fq6Mul(c, a, b *fq6) {}
func fq6Sqr(c, a, b *fq6) {}
