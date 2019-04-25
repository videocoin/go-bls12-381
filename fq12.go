package bls12

import "math/big"

// fq12 implements the field of size q¹² as a quadratic extension of fq6 where
// y = v.
type fq12 struct {
	c0, c1 fq6
}

// SetOne sets z to 1 and returns z.
func (z *fq12) SetOne() *fq12 {
	z.c0.SetOne()
	z.c1.SetZero()
	return z
}

// Set sets z to x and returns z.
func (z *fq12) Set(x *fq12) *fq12 {
	z.c0 = x.c0
	z.c1 = x.c1
	return x
}

// Equal reports whether x is equal to y.
func (x *fq12) Equal(y *fq12) bool {
	return x.c0 == y.c0 && x.c1 == y.c1
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
	elem := new(fq12)
	v0 := new(fq6).Mul(&x.c0, &y.c0)
	v1 := new(fq6).Mul(&x.c1, &y.c1)
	elem.c0.MulQuadraticNonResidue(v1).Add(&elem.c0, v1)
	t0 := new(fq6)
	elem.c1.Add(&x.c0, &x.c1)
	elem.c1.Mul(&elem.c1, t0.Add(&y.c0, &y.c1)).Sub(&elem.c1, t0.Add(v0, v1))

	return z.Set(elem)
}

// SparseMult sets z to the product of x with a0, a1, a2 and returns z.
// SparseMult utilizes the sparness property to avoid full Fq12 arithmetic.
// See https://eprint.iacr.org/2012/408.pdf - algorithm 5.
func (z *fq12) SparseMult(x *fq12, a0 *fq2, a1 *fq2, a2 *fq2) *fq12 {
	a := new(fq6)
	a.c0.Mul(a0, &x.c0.c0)
	a.c1.Mul(a0, &x.c0.c1)
	a.c2.Mul(a0, &x.c0.c2)
	b := new(fq6).SparseMul(&x.c1, a1, a2)
	c := new(fq6)
	c.c0.Add(a0, a1)
	c.c1 = *a2
	d := new(fq6).Add(&x.c0, &x.c1)
	e := new(fq6).SparseMul(d, &c.c0, &c.c1)
	f := new(fq6).Sub(e, new(fq6).Add(a, b))
	g := new(fq6).MulQuadraticNonResidue(b)
	h := new(fq6).Add(a, g)
	z.c0 = *h
	z.c1 = *f

	return z
}

// Sqr sets z to the product x*x and returns z.
// Sqr utilizes complex squaring - See https://eprint.iacr.org/2006/471.pdf.
// TODO
// https://books.google.pt/books?id=6utNZkI-oGkC&pg=PA193&lpg=PA193&dq=karatsuba+or+complex+for+quadratic+extensions+squaring&source=bl&ots=twAZiIiLPG&sig=ACfU3U25YBNFAhdb0j8iUTVwekou6oMByQ&hl=en&sa=X&ved=2ahUKEwjJ_tT2lNfhAhVSUxoKHQZBBfgQ6AEwAnoECAgQAQ#v=onepage&q=karatsuba%20or%20complex%20for%20quadratic%20extensions%20squaring&f=false
func (z *fq12) Sqr(x *fq12) *fq12 {
	// v0 = a0a1
	v0 := new(fq6).Mul(&x.c0, &x.c1)
	// c0 = (a0 + a1)(a0 + γa1) − v0 − γv0
	c, t0 := new(fq12), new(fq6)
	c.c0.Add(&x.c0, &x.c1).Mul(&c.c0, t0.MulQuadraticNonResidue(&x.c1).Add(t0, &x.c0)).Sub(&c.c0, t0.Add(v0, t0.MulQuadraticNonResidue(v0)))
	// c1 = 2v0
	c.c1.Add(v0, v0)

	return z.Set(c)
}

// Inv sets z to 1/x and returns z.
// See "Implementing cryptographic pairings", M. Scott - section 3.2.
func (z *fq12) Inv(x *fq12) *fq12 {
	t0, t1 := new(fq6).Sqr(&x.c0), new(fq6).Sqr(&x.c1)
	t1.MulQuadraticNonResidue(t1).Sub(t0, t1).Inv(t1)
	z.c0.Mul(&x.c0, t1)
	z.c1.Mul(&x.c1, t1)
	return z
}

// Exp sets z=x**y and returns z.
func (z *fq12) Exp(x *fq12, y *big.Int) *fq12 {
	b := *x
	ret := new(fq12).SetOne()
	for i := y.BitLen() - 1; i >= 0; i-- {
		if y.Bit(i) == 1 {
			ret.Mul(ret, &b)
		}
		b.Sqr(&b)
	}

	return z.Set(ret)
}

// Frobenius sets z to the pth-power Frobenius of x and returns z.
func (z *fq12) Frobenius(x *fq12, power uint64) *fq12 {
	// TODO
	if power == 6 {
		return z.Conjugate(x)
	}

	return &fq12{}
}
