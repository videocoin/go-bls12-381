package bls12

// fq6 is an element of Fq6 = Fq²[Y]/(Y³ − γ), where γ
// is a quadratic non-residue in Fq and γ = √β is a
// cubic non-residue in Fq² with a value of X + 1.
// See https://eprint.iacr.org/2006/471.pdf - "6.2 Cubic over quadratic"
type fq6 struct {
	c0, c1, c2 fq2
}

func (x *fq6) SetOne() *fq6 {
	x.c0.SetOne()
	x.c1.SetZero()
	x.c2.SetZero()
	return x
}

func (x *fq6) SetZero() *fq6 {
	x.c0.SetZero()
	x.c1.SetZero()
	x.c2.SetZero()
	return x
}

func (z *fq6) Sub(x, y *fq6) *fq6 {
	z.c0.Sub(&x.c0, &y.c0)
	z.c1.Sub(&x.c1, &y.c1)
	z.c2.Sub(&x.c2, &y.c2)
	return z
}

func (z *fq6) Add(x, y *fq6) *fq6 {
	z.c0.Add(&x.c0, &y.c0)
	z.c1.Add(&x.c1, &y.c1)
	z.c2.Add(&x.c2, &y.c2)
	return z
}

func (z *fq6) Dbl(x *fq6) *fq6 {
	return z.Add(x, x)
}

// Cubic extensions - Karatsuba method
// See https://eprint.iacr.org/2006/471.pdf - Page 6,7
// fixme: B value
func (z *fq6) Mul(x, y *fq6) *fq6 {
	v0, v1, v2 := new(fq2), new(fq2), new(fq2)
	v0.Mul(&x.c0, &y.c0)
	v1.Mul(&x.c1, &y.c1)
	v2.Mul(&x.c2, &y.c2)

	c0, c1, c2, t0, t1 := new(fq2), new(fq2), new(fq2), new(fq2), new(fq2)
	t0.Add(&x.c1, &x.c2)
	t1.Add(&y.c1, &y.c2)
	t0.Mul(t0, t1)
	t1.Add(v1, v2)
	t0.Sub(t0, t1)
	//fq2Mul(c0, )
	c0.Add(c0, v0)

	t0.Add(&x.c0, &x.c1)
	t1.Add(&y.c0, &y.c1)
	t0.Mul(t0, t1)
	t1.Add(v0, v1)
	t0.Sub(t0, t1)
	//fq2Mul()
	c1.Add(t0, t1)

	t0.Add(&x.c0, &x.c2)
	t1.Add(&y.c0, &y.c2)
	t0.Mul(t0, t1)
	c2.Add(t0, v1)
	t1.Add(v0, v2)
	c2.Sub(c2, t1)

	z.c0, z.c1, z.c2 = *c0, *c1, *c2

	return z
}
