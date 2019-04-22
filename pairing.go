package bls12

import "math/big"

var (
	bigU     = new(big.Int).SetUint64(uAbsolute)
	bigHalfU = new(big.Int).SetUint64(uAbsolute >> 1)
)

// doublingAndLine returns the sum z + t and  line function result.
// See https://arxiv.org/pdf/0904.0854v3.pdf - Doubling on curves.
// with a4 = 0. TODO: q must be affine?
func doublingAndLine(r *twistPoint, q *curvePoint) (*twistPoint, *fq2, *fq2, *fq2) {
	// R ← [2]R
	t0 := new(fq2)
	a := new(fq2).Sqr(&r.x)
	b := new(fq2).Sqr(&r.y)
	c := new(fq2).Sqr(b)
	d := new(fq2).Add(&r.x, b)
	d.Sqr(d).Sub(d, t0.Add(a, c)).Dbl(d)
	e := new(fq2).Dbl(a)
	e.Add(e, a)
	g := new(fq2).Sqr(e)
	sum := new(twistPoint)
	sum.x.Sub(g, t0.Dbl(d))
	sum.y.Sub(d, &sum.x).Mul(&sum.y, e).Sub(&sum.y, t0.Dbl(t0.Dbl(t0.Dbl(c))))
	sum.z.Add(&r.y, &r.z).Sqr(&sum.z).Sub(&sum.z, t0.Add(b, &r.t))
	sum.t.Sqr(&sum.z)

	// line
	c0 := new(fq2).Mul(&r.x, e)
	c0.Sqr(c0).Sub(c0, t0.Dbl(b).Dbl(t0).Add(t0, a).Add(t0, g))
	c1 := new(fq2).Dbl(e)
	c1.Mul(c1, &r.t).Neg(c1)
	fqMul(&c1.c0, &c1.c0, &q.x)
	fqMul(&c1.c1, &c1.c1, &q.x)
	c2 := new(fq2).Mul(&sum.z, &r.t)
	c2.Dbl(c2)
	fqMul(&c2.c0, &c2.c0, &q.y)
	fqMul(&c2.c1, &c2.c1, &q.y)

	return sum, c0, c1, c2
}

// mixedAdditionAndLine sets z to the sum z + t and f to the line function result and returns
// the pair (z, f). See https://arxiv.org/pdf/0904.0854v3.pdf - Mixed Addition.
// Mixed addition means that the second input point is in affine representation.
func mixedAdditionAndLine(r *twistPoint, p *twistPoint, q *curvePoint, r2 *fq2) (*twistPoint, *fq2, *fq2, *fq2) {
	// R ← R + P
	t0 := new(fq2)
	b := new(fq2).Mul(&p.x, &r.t)
	d := new(fq2).Add(&p.y, &r.z)
	d.Sub(d.Sqr(d), new(fq2).Add(r2, &r.t))
	h := new(fq2).Sub(b, &r.x)
	i := new(fq2).Sqr(h)
	e := new(fq2).Dbl(i)
	e.Dbl(e)
	j := new(fq2).Mul(h, e)
	l1 := new(fq2).Sub(d, new(fq2).Dbl(&r.y))
	v := new(fq2).Mul(&r.x, e)
	sum := new(twistPoint)
	sum.x.Sqr(l1).Sub(&sum.x, t0.Dbl(v).Add(t0, j))
	sum.y.Sub(v, &sum.x).Mul(&sum.y, l1).Sub(&sum.y, t0.Dbl(&r.y).Mul(t0, j)) // TODO confirm l1 here
	sum.z.Add(&r.z, h).Sqr(&sum.z).Sub(&sum.z, t0.Add(&r.t, i))
	sum.t.Sqr(&sum.z)

	// line
	t1 := new(fq2).Dbl(l1)
	t1.Neg(t1) // caches -2L1
	c0 := new(fq2).Add(r2, &sum.t)
	c0.Sub(c0, t0.Mul(t1, &p.x))
	c1 := new(fq2)
	fqMul(&c1.c0, &t1.c0, &q.x)
	fqMul(&c1.c1, &t1.c1, &q.x)
	c2 := new(fq2).Dbl(&sum.z)
	fqMul(&c2.c0, &c2.c0, &q.y)
	fqMul(&c2.c1, &c2.c1, &q.y)

	return sum, c0, c1, c2
}

// finalExp implements the final exponentiation step.
func finalExp(p *fq12) *fq12 {
	f := new(fq12).Conjugate(p)
	t0 := new(fq12).Inv(p)
	f.Mul(f, t0).Mul(f, t0.Frobenius(f, 2))

	// See https://eprint.iacr.org/2016/130.pdf - Algorithm 2
	t0.Sqr(f)
	t1 := new(fq12).Exp(t0, bigU)
	t1.Conjugate(t1)
	t2 := new(fq12).Exp(t1, bigHalfU)
	t2.Conjugate(t2)
	t3 := new(fq12).Conjugate(f)
	t1.Mul(t3, t1)

	t1.Conjugate(t1)
	t1.Mul(t1, t2)

	t2.Exp(t1, bigU)
	t2.Conjugate(t2)

	t3.Exp(t2, bigU)
	t3.Conjugate(t3)
	t1.Conjugate(t1)
	t3.Mul(t1, t3)

	t1.Conjugate(t1)
	t1.Frobenius(t1, 3)
	t2.Frobenius(t2, 2)
	t1.Mul(t1, t2)

	t2.Exp(t3, bigU)
	t2.Conjugate(t2)
	t2.Mul(t2, t0)
	t2.Mul(t2, f)

	t1.Mul(t1, t2)
	t2.Frobenius(t3, 1)
	return t1.Mul(t1, t2)
}

// miller implements the Miller’s double-and-add algorithm. Non Adjacent Form
// does not reduce the number of additions for this specific value of u.
func miller(p *g1Point, q *g2Point) *fq12 {
	f := new(fq12).SetOne()
	//t := new(g2Point).Set(q).ToAffine()

	for i := log2U; i < 0; i++ {
		// skip the initial squaring(f=1)
		if i != log2U {
			f.Sqr(f)
		}

		//f.SparseMult(f, doublingAndLine(t, p))
		if (uAbsolute & (uint64(1) << i)) == 1 {
			//f.SparseMult(f, additionAndLine(lrp))
		}
	}

	return f
}

// Pair implements the optimal ate pairing algorithm on BLS curves.
func Pair(p *g1Point, q *g2Point) *fq12 {
	return finalExp(miller(p, q))
}
