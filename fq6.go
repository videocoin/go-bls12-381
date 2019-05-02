package bls12

// fq6 is an element of Fq6 = Fq²[Y]/(Y³ − γ), where γ
// is a quadratic non-residue in Fq and γ = √β is a
// cubic non-residue in Fq² with a value of X + 1.
// See https://eprint.iacr.org/2006/471.pdf - "6.2 Cubic over quadratic"
type fq6 struct {
	c0, c1, c2 fq2
}

// SetOne sets z to 0 and returns z.
func (z *fq6) SetZero() *fq6 {
	z.c0.SetZero()
	z.c1.SetZero()
	z.c2.SetZero()
	return z
}

// SetOne sets z to 1 and returns z.
func (z *fq6) SetOne() *fq6 {
	z.c0.SetOne()
	z.c1.SetZero()
	z.c2.SetZero()
	return z
}

// Set sets z to x and returns z.
func (z *fq6) Set(x *fq6) *fq6 {
	z.c0.Set(&x.c0)
	z.c1.Set(&x.c1)
	z.c2.Set(&x.c2)
	return z
}

// Neg sets z to -x and returns z.
func (z *fq6) Neg(x *fq6) *fq6 {
	z.c0.Neg(&x.c0)
	z.c1.Neg(&x.c1)
	z.c2.Neg(&x.c2)
	return z
}

// Add sets z to the sum x+y and returns z.
func (z *fq6) Add(x, y *fq6) *fq6 {
	z.c0.Add(&x.c0, &y.c0)
	z.c1.Add(&x.c1, &y.c1)
	z.c2.Add(&x.c2, &y.c2)
	return z
}

// Sub sets z to the difference x-y and returns z.
func (z *fq6) Sub(x, y *fq6) *fq6 {
	z.c0.Sub(&x.c0, &y.c0)
	z.c1.Sub(&x.c1, &y.c1)
	z.c2.Sub(&x.c2, &y.c2)
	return z
}

// Mul sets z to the product x*y and returns z.
// Mul utilizes Karatsuba's method.
// See https://eprint.iacr.org/2006/471.pdf page 6-7.
// TODO
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
	c0.MulXi(t0).Add(c0, v0)

	t1.Add(&y.c0, &y.c1)
	t0.Add(&x.c0, &x.c1).Mul(t0, t1)
	t1.Add(v0, v1)
	t0.Sub(t0, t1)
	c1.MulXi(v2).Add(c1, t0)

	t0.Add(&x.c0, &x.c2)
	t1.Add(&y.c0, &y.c2)
	t0.Mul(t0, t1)
	c2.Add(t0, v1)
	t1.Add(v0, v2)
	c2.Sub(c2, t1)

	z.c0, z.c1, z.c2 = *c0, *c1, *c2

	return z
}

// SparseMult sets z to the product of x with a0, a1, a2 and returns z.
// SparseMult utilizes the sparness property to avoid full Fq6 arithmetic.
// See https://eprint.iacr.org/2012/408.pdf - Algorithm 6.
func (z *fq6) SparseMul(x *fq6, a0 *fq2, a1 *fq2) *fq6 {
	ret, t0 := new(fq6), new(fq2)

	a := new(fq2).Mul(a0, &x.c0)
	b := new(fq2).Mul(a1, &x.c1)
	e := new(fq2).Add(a0, a1)
	e.Mul(e, t0.Add(&x.c0, &x.c1))
	ret.c0.Add(a, t0.MulXi(&x.c2).Mul(t0, a1))
	ret.c1.Sub(e, t0.Add(a, b))
	ret.c2.Mul(a0, &x.c2).Add(&z.c2, b)

	return z.Set(ret)
}

// MulQuadraticNonResidue sets z to the product γX and returns z.
func (z *fq6) MulQuadraticNonResidue(x *fq6) *fq6 {
	// γ = v
	// X = a0 + a1v + a2v^2
	// γX = a0v + a1v^2 + a2ξ
	ret := new(fq6)
	ret.c0.MulXi(&x.c2)
	ret.c1.Set(&x.c0)
	ret.c2.Set(&x.c1)

	return z.Set(ret)
}

// Sqr sets z to the product x*x and returns z.
// Sqr utilizes the CH-SQR3x method.
// See https://eprint.iacr.org/2006/471.pdf page 9.
func (z *fq6) Sqr(x *fq6) *fq6 {
	s0 := new(fq2).Sqr(&x.c0)
	s1 := new(fq2).Add(&x.c0, &x.c1)
	s1.Add(s1, &x.c2).Sqr(s1)
	s2 := new(fq2).Add(&x.c0, &x.c2)
	s2.Sub(s2, &x.c1).Sqr(s2)
	s3 := new(fq2).Mul(&x.c1, &x.c2)
	s3.Add(s3, s3)
	s4 := new(fq2).Sqr(&x.c2)

	t0, t1 := new(fq2), new(fq2)
	z.c0.Add(s0, s0).Add(&z.c0, t0.MulXi(s3).Add(t0, t0))
	z.c1.Sub(s1, t0.Add(s2, t0.Add(s3, s3))).Add(&z.c1, t0.MulXi(s4).Add(t0, t0))
	z.c2.Add(s1, s2).Sub(&z.c2, t0.Add(t0.Add(s0, s0), t1.Add(s4, s4)))

	return z
}

func (z *fq6) Inv(x *fq6) *fq6 {
	// TODO
	return &fq6{}
}

func (z *fq6) Frobenius(x *fq6, power uint64) *fq6 {
	ret := new(fq6)
	ret.c0.Frobenius(&x.c0, power)
	ret.c1.Frobenius(&x.c1, power)
	ret.c1.Mul(&ret.c1, frobeniusCoeff6c1[power%6])
	ret.c2.Frobenius(&x.c2, power)
	ret.c2.Mul(&ret.c2, frobeniusCoeff6c2[power%6])
	return z.Set(ret)
}
