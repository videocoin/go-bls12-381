package bls12

import "math/big"

// fq12 implements the field of size q¹² as a quadratic extension of fq6 where
// y = v. See https://eprint.iacr.org/2006/471.pdf for arithmetic.
type fq12 struct {
	c0, c1 fq6
}

// Set sets z to x and returns z.
func (z *fq12) Set(x *fq12) *fq12 {
	z.c0.Set(&x.c0)
	z.c1.Set(&x.c1)
	return z
}

// SetOne sets z to 1 and returns z.
func (z *fq12) SetOne() *fq12 {
	z.c0.SetOne()
	z.c1.SetZero()
	return z
}

// Equal reports whether x is equal to y.
func (x *fq12) Equal(y *fq12) bool {
	return *x == *y
}

// Conjugate sets z to the conjugate of x and returns z.
func (z *fq12) Conjugate(x *fq12) *fq12 {
	z.c0.Set(&x.c0)
	z.c1.Neg(&x.c1)
	return z
}

// Add sets z to the sum x+y and returns z.
func (z *fq12) Add(x, y *fq12) *fq12 {
	z.c0.Add(&x.c0, &y.c0)
	z.c1.Add(&x.c1, &y.c1)

	return z
}

// Mul sets z to the product x*y and returns z.
// Mul utilizes Karatsuba's method.
func (z *fq12) Mul(x, y *fq12) *fq12 {
	ret, t0, v1 := new(fq12), new(fq6), new(fq6).Mul(&x.c1, &y.c1)
	ret.c0.MulQuadraticNonResidue(v1).Add(&ret.c0, v1)
	ret.c1.Add(&x.c0, &x.c1).Mul(&ret.c1, t0.Add(&y.c0, &y.c1)).Sub(&ret.c1, t0.Add(t0.Mul(&x.c0, &y.c0), v1))

	return z.Set(ret)
}

// SparseMult sets z to the product of x with a0, a1, a2 and returns z.
// SparseMult utilizes the sparness property to avoid full fq12 arithmetic.
// See https://eprint.iacr.org/2012/408.pdf - algorithm 5.
func (z *fq12) SparseMult(x *fq12, a0 *fq2, a1 *fq2, a2 *fq2) *fq12 {
	a, t0 := new(fq6), new(fq6)
	a.c0.Mul(a0, &x.c0.c0)
	a.c1.Mul(a0, &x.c0.c1)
	a.c2.Mul(a0, &x.c0.c2)
	b := new(fq6).SparseMul(&x.c1, a1, a2)
	c := new(fq6)
	c.c0.Add(a0, a1)
	c.c1 = *a2
	e := new(fq6).SparseMul(t0.Add(&x.c0, &x.c1), &c.c0, &c.c1) // d, e
	z.c1.Sub(e, t0.Add(a, b))                                   // f
	z.c0.Add(a, t0.MulQuadraticNonResidue(b))                   // g, h

	return z
}

// Sqr sets z to the product x*x and returns z.
// Sqr utilizes complex squaring.
func (z *fq12) Sqr(x *fq12) *fq12 {
	ret, tmp, v0 := new(fq12), new(fq6), new(fq6).Mul(&x.c0, &x.c1)
	ret.c0.Add(&x.c0, &x.c1).Mul(&ret.c0, tmp.MulQuadraticNonResidue(&x.c1).Add(tmp, &x.c0)).Sub(&ret.c0, tmp.Add(v0, tmp.MulQuadraticNonResidue(v0)))
	ret.c1.Add(v0, v0)

	return z.Set(ret)
}

// Inv sets z to 1/x and returns z.
// See "Implementing cryptographic pairings", M. Scott - section 3.2.
func (z *fq12) Inv(x *fq12) *fq12 {
	tmp0, tmp1 := new(fq6).Sqr(&x.c0), new(fq6).Sqr(&x.c1)
	tmp1.MulQuadraticNonResidue(tmp1).Sub(tmp0, tmp1).Inv(tmp1)
	z.c0.Mul(&x.c0, tmp1)
	z.c1.Mul(&x.c1, tmp1)

	return z
}

// Exp sets z=x**y and returns z.
func (z *fq12) Exp(x *fq12, y *big.Int) *fq12 {
	ret, base := new(fq12).SetOne(), *x
	for i := y.BitLen() - 1; i >= 0; i-- {
		if y.Bit(i) == 1 {
			ret.Mul(ret, &base)
		}
		base.Sqr(&base)
	}

	return z.Set(ret)
}

// Frobenius sets z to the pth-power Frobenius of x and returns z.
func (z *fq12) Frobenius(x *fq12, power uint64) *fq12 {
	ret := new(fq12)
	ret.c0.Frobenius(&x.c0, power)
	ret.c1.Frobenius(&x.c1, power)
	ret.c1.c0.Mul(&ret.c1.c0, frob12c1[power%12])
	ret.c1.c1.Mul(&ret.c1.c1, frob12c1[power%12])
	ret.c1.c2.Mul(&ret.c1.c2, frob12c1[power%12])

	return z.Set(ret)
}
