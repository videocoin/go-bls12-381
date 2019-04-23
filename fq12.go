package bls12

import "math/big"

// fq12 implements the field of size q¹² as a quadratic extension of fq6
// where γ = v.
type fq12 struct {
	c0, c1 fq6
}

func (z *fq12) SetOne() *fq12 {
	z.c0.SetOne()
	z.c0.SetZero()
	return z
}

func (z *fq12) Add(x, y *fq12) *fq12 {
	z.c0.Add(&x.c0, &y.c0)
	z.c1.Add(&x.c1, &y.c1)
	return z
}

// https://books.google.pt/books?id=6utNZkI-oGkC&pg=PA193&lpg=PA193&dq=karatsuba+or+complex+for+quadratic+extensions+squaring&source=bl&ots=twAZiIiLPG&sig=ACfU3U25YBNFAhdb0j8iUTVwekou6oMByQ&hl=en&sa=X&ved=2ahUKEwjJ_tT2lNfhAhVSUxoKHQZBBfgQ6AEwAnoECAgQAQ#v=onepage&q=karatsuba%20or%20complex%20for%20quadratic%20extensions%20squaring&f=false
// See https://eprint.iacr.org/2006/471.pdf - Quadratic Extensions; Complex Squaring
func (z *fq12) Sqr(x *fq12) *fq12 {
	// v0 = a0a1
	v0 := new(fq6).Mul(&x.c0, &x.c1)
	// c0 = (a0 + a1)(a0 + γa1) − v0 − γv0
	c, t0 := new(fq12), new(fq6)
	c.c0.Add(&x.c0, &x.c1).Mul(&c.c0, t0.MulQNR(&x.c1).Add(t0, &x.c0)).Sub(&c.c0, t0.Add(v0, t0.MulQNR(v0)))
	// c1 = 2v0
	c.c1.Dbl(v0)

	return z.Set(c)
}

func (x *fq12) Equal(y *fq12) bool {
	return x.c0 == y.c0 && x.c1 == y.c1
}

func (x *fq12) Set(y *fq12) *fq12 {
	x.c0 = y.c0
	x.c1 = y.c1
	return x
}

// karatsuba method
func (z *fq12) Mul(x, y *fq12) *fq12 {
	elem := new(fq12)
	v0 := new(fq6).Mul(&x.c0, &y.c0)
	v1 := new(fq6).Mul(&x.c1, &y.c1)
	elem.c0.MulQNR(v1).Add(&elem.c0, v1)
	t0 := new(fq6)
	elem.c1.Add(&x.c0, &x.c1)
	elem.c1.Mul(&elem.c1, t0.Add(&y.c0, &y.c1)).Sub(&elem.c1, t0.Add(v0, v1))

	return z.Set(elem)
}

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
	g := new(fq6).MulQNR(b)
	h := new(fq6).Add(a, g)
	z.c0 = *h
	z.c1 = *f

	return z
}

// Inv
// See "Implementing cryptographic pairings", M. Scott - section 3.2.
func (z *fq12) Inv(x *fq12) *fq12 {
	t0, t1 := new(fq6).Sqr(&x.c0), new(fq6).Sqr(&x.c1)
	t1.MulQNR(t1).Sub(t0, t1).Inv(t1)
	z.c0.Mul(&x.c0, t1)
	z.c1.Mul(&x.c1, t1)
	return z
}

func (z *fq12) Exp(x *fq12, y *big.Int) *fq12 {
	return &fq12{}
}

/*
func (z *fq12) Exp(x *fq12, y []uint64) *fq12 {
	/*
		b := *x
		ret := new(fq12).SetOne()
		for _, word := range y {
			for j := uint(0); j < wordSize; j++ {
				if (word & (1 << j)) != 0 {
					fqMul(ret, ret, &b)
				}
				fqSqr(&b, &b)
			}
		}

	return &fq12{}
}
*/

func (z *fq12) Frobenius(x *fq12, power uint64) *fq12 {
	// TODO
	return &fq12{}
}

func (z *fq12) Conjugate(x *fq12) *fq12 {
	z.c0.Set(&x.c0)
	z.c1.Neg(&x.c1)

	return z
}
