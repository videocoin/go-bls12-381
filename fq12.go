package bls12

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

func (z *fq12) Mul(x, y *fq12) *fq12 {
	// TODO
	return &fq12{}
}

// SparseMult implements the 8-sparse multiplication.
// See https://eprint.iacr.org/2017/1174.pdf - Algorithm 2.
func (z *fq12) SparseMult(x *fq12, y *fq12) *fq12 {
	c := new(fq12)
	// c4 ← a0 × b4
	c.c1.c1.Mul(&x.c0.c0, &y.c1.c1)
	// t1 ← a1 × b5,
	t1 := new(fq2).Mul(&x.c0.c1, &y.c1.c2)
	// t2 ← a0 + a1,
	t2 := new(fq2).Add(&x.c0.c0, &x.c0.c1)
	// S0 ← b4 + b5
	s0 := new(fq2).Add(&y.c1.c1, &y.c1.c2)
	// c5 ← t2 × S0 − (c4 + t1)
	c.c1.c2.Mul(t2, s0).Sub(&c.c1.c2, new(fq2).Add(&c.c1.c1, t1))
	// t0 ← a2 × b4
	t0 := new(fq2).Mul(&x.c1.c2, &y.c1.c1)
	// t0 ← t0 + t1
	t0.Add(t0, t1)
	// t0 ← a3 × b4
	t0.Mul(&x.c1.c0, &y.c1.c1)
	// t1 ← a4 × b5
	t1.Mul(&x.c1.c1, &y.c1.c2)
	// t2 ← a3 + a4
	t2.Add(&x.c1.c0, &y.c1.c1)
	// c2 ← t0
	c.c0.c2.Add(t0, t1)
	// c ← c + a
	c.Add(c, x)

	return z.Set(c)
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
