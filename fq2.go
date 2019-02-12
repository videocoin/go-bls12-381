package bls12

// TODO (rgeraldes) - desc
// fq2 is an element (c0 + c1 * u)
type fq2 struct {
	c0, c1 fq
}

func fq2Add(c, a, b *fq2) {}
func fq2Mul(c, a, b *fq2) {}
func fq2Sub(c, a, b *fq2) {}
func fq2Sqr(c, a *fq2)    {}
func fq2Dbl(c, a *fq2)    {}
