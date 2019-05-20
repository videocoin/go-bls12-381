package bls12

// fq6 is an element of Fq6 = Fq²[Y]/(Y³ − γ), where γ is a quadratic
// non-residue in Fq and γ = √β is a cubic non-residue in Fq² with a value of
// X + 1. See https://eprint.iacr.org/2006/471.pdf for arithmetic.
type fq6 struct {
	c0, c1, c2 fq2
}

// Set sets z to x and returns z.
func (z *fq6) Set(x *fq6) *fq6 {
	z.c0.Set(&x.c0)
	z.c1.Set(&x.c1)
	z.c2.Set(&x.c2)
	return z
}

// SetZero sets z to 0 and returns z.
func (z *fq6) SetZero() *fq6 {
	*z = fq6{}
	return z
}

// SetOne sets z to 1 and returns z.
func (z *fq6) SetOne() *fq6 {
	z.c0.SetOne()
	z.c1.SetZero()
	z.c2.SetZero()
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
func (z *fq6) Mul(x, y *fq6) *fq6 {
	ret, t0, t1 := new(fq6), new(fq2), new(fq2)
	v0 := new(fq2).Mul(&x.c0, &y.c0)
	v1 := new(fq2).Mul(&x.c1, &y.c1)
	v2 := new(fq2).Mul(&x.c2, &y.c2)
	// c0
	t1.Add(&y.c1, &y.c2)
	t0.Add(&x.c1, &x.c2).Mul(t0, t1)
	t1.Add(v1, v2)
	t0.Sub(t0, t1)
	ret.c0.MulXi(t0).Add(&ret.c0, v0)
	// c1
	t1.Add(&y.c0, &y.c1)
	t0.Add(&x.c0, &x.c1).Mul(t0, t1)
	t1.Add(v0, v1)
	t0.Sub(t0, t1)
	ret.c1.MulXi(v2).Add(&ret.c1, t0)
	// c2
	t1.Add(&y.c0, &y.c2)
	t0.Add(&x.c0, &x.c2).Mul(t0, t1)
	t1.Add(v0, v2)
	ret.c2.Add(t0, v1).Sub(&ret.c2, t1)

	return z.Set(ret)
}

/*
// SparseMult sets z to the product of x with a0, a1, a2 and returns z.
// SparseMult utilizes the sparness property to avoid full fq6 arithmetic.
// See https://eprint.iacr.org/2012/408.pdf - Algorithm 6.
func (z *fq6) SparseMul(x *fq6, a0 *fq2, a1 *fq2) *fq6 {
	ret, t0 := new(fq6), new(fq2)
	a := new(fq2).Mul(a0, &x.c0)
	b := new(fq2).Mul(a1, &x.c1)
	e := new(fq2).Add(a0, a1)
	e.Mul(e, t0.Add(&x.c0, &x.c1))
	ret.c0.Add(a, t0.MulXi(&x.c2).Mul(t0, a1)) // d
	ret.c1.Sub(e, t0.Add(a, b))                // g
	ret.c2.Mul(a0, &x.c2).Add(&z.c2, b)        // i

	return z.Set(ret)
}*/

// SparseMult01 sets z to the product of x with c0, c1 and returns z.
// See https://github.com/zkcrypto/pairing/blob/master/src/bls12_381/fq6.rs#L68.
func (z *fq6) SparseMul01(x *fq6, c0 *fq2, c1 *fq2) *fq6 {
	aa := new(fq2).Mul(&x.c0, c0)
	bb := new(fq2).Mul(&x.c1, c1)
	ret, t0 := new(fq6), new(fq2)
	ret.c0.Mul(c1, t0.Add(&x.c1, &x.c2)).Sub(&ret.c0, bb).MulXi(&ret.c0).Add(&ret.c0, aa)
	ret.c2.Mul(c0, t0.Add(&x.c0, &x.c2)).Sub(&ret.c2, aa).Add(&ret.c2, bb)
	ret.c1.Add(c0, c1).Mul(&ret.c1, t0.Add(&x.c0, &x.c1)).Sub(&ret.c1, t0.Add(aa, bb))

	return z.Set(ret)
}

// SparseMult1 sets z to the product of x with c0, c1 and returns z.
// See https://github.com/zkcrypto/pairing/blob/master/src/bls12_381/fq6.rs#L40.
func (z *fq6) SparseMul1(x *fq6, c1 *fq2) *fq6 {
	ret, t0 := new(fq6), new(fq2)
	ret.c2.Mul(&x.c1, c1)
	ret.c0.Mul(c1, t0.Add(&x.c1, &x.c2)).Sub(&ret.c0, &ret.c2).MulXi(&ret.c0)
	ret.c1.Mul(c1, t0.Add(&x.c0, &x.c1)).Sub(&ret.c1, &ret.c2)

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

/*
// TODO
// Sqr sets z to the product x*x and returns z.
// Sqr utilizes the CH-SQR3x method (c = 2a^2).
func (z *fq6) Sqr(x *fq6) *fq6 {
	t0, t1 := new(fq2), new(fq2)
	s0 := new(fq2).Sqr(&x.c0)
	s1 := new(fq2).Add(&x.c0, &x.c1)
	s1.Add(s1, &x.c2).Sqr(s1)
	s2 := new(fq2).Add(&x.c0, &x.c2)
	s2.Sub(s2, &x.c1).Sqr(s2)
	s3 := new(fq2).Mul(&x.c1, &x.c2)
	s3.Add(s3, s3)
	s4 := new(fq2).Sqr(&x.c2)
	// c0
	z.c0.Add(s0, s0).Add(&z.c0, t0.MulXi(s3).Add(t0, t0))
	// c1
	z.c1.Sub(s1, t0.Add(s2, t0.Add(s3, s3))).Add(&z.c1, t0.MulXi(s4).Add(t0, t0))
	// c2
	z.c2.Add(s1, s2).Sub(&z.c2, t0.Add(t0.Add(s0, s0), t1.Add(s4, s4)))

	return z
}
*/

// Sqr sets z to the product x*x and returns z.
// Sqr utilizes Karatsuba's method (c = 2a^2).
// TODO fastest method?
func (z *fq6) Sqr(x *fq6) *fq6 {
	v0 := new(fq2).Sqr(&x.c0)
	v1 := new(fq2).Sqr(&x.c1)
	v2 := new(fq2).Sqr(&x.c2)
	ret, t0 := new(fq6), new(fq2)
	ret.c0.Add(&x.c1, &x.c2).Sqr(&ret.c0).Sub(&ret.c0, t0.Add(v1, v2)).MulXi(&ret.c0).Add(&ret.c0, v0)
	ret.c1.Add(&x.c0, &x.c1).Sqr(&ret.c1).Sub(&ret.c1, t0.Add(v0, v1)).Add(&ret.c1, t0.MulXi(v2))
	ret.c2.Add(&x.c0, &x.c2).Sqr(&ret.c2).Sub(&ret.c2, t0.Add(v0, v2)).Add(&ret.c2, v1)

	return z.Set(ret)
}

// Inv sets z to 1/x and returns z.
func (z *fq6) Inv(x *fq6) *fq6 {
	// v0 = x0^2 - E * x1 * x2
	// v1 = E * x2^2 - x0 * x1
	// v2 = x1^2 - x0 * x2
	// z0 = x0 * v0
	// z1 = x1 * v2 * E
	// z2 = x2 * v1 * E
	// t0 = 1/(z0 + z1 + z2)
	// z0 = v0 * t0
	// z1 = v1 * t0
	// z2 = v2 * t0
	ret, t0 := new(fq6), new(fq2)
	v0 := new(fq2).Mul(&x.c1, &x.c2)
	v0.MulXi(v0).Sub(t0.Sqr(&x.c0), v0)
	v1 := new(fq2).Sqr(&x.c2)
	v1.MulXi(v1).Sub(v1, t0.Mul(&x.c0, &x.c1))
	v2 := new(fq2).Sqr(&x.c1)
	v2.Sub(v2, t0.Mul(&x.c0, &x.c2))
	ret.c0.Mul(&x.c0, v0)
	ret.c1.Mul(&x.c1, v2).MulXi(&ret.c1)
	ret.c2.Mul(&x.c2, v1).MulXi(&ret.c2)
	t0.Add(&ret.c0, &ret.c1).Add(t0, &ret.c2).Inv(t0)
	ret.c0.Mul(v0, t0)
	ret.c1.Mul(v1, t0)
	ret.c2.Mul(v2, t0)

	return z.Set(ret)
}

// Frobenius sets z to frobenius x for a certain power and returns z.
func (z *fq6) Frobenius(x *fq6, power uint64) *fq6 {
	ret := new(fq6)
	ret.c0.Frobenius(&x.c0, power)
	ret.c1.Frobenius(&x.c1, power).Mul(&ret.c1, frobFq6C1[power%6])
	ret.c2.Frobenius(&x.c2, power).Mul(&ret.c2, frobFq6C2[power%6])

	return z.Set(ret)
}
