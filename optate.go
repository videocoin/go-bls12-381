package bls12

import (
	"math/big"
)

var (
	bigU     = new(big.Int).SetUint64(15132376222941642752)
	bigHalfU = new(big.Int).Rsh(bigU, 1)
	uArr     = []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 1, 1}
	uArrLen  = len(uArr)
)

// doublingAndLine returns the sum r + r and the line function result.
// See https://arxiv.org/pdf/0904.0854v3.pdf - Doubling on curves with a4 = 0.
func doublingAndLine(r *twistPoint, q *curvePoint) (*twistPoint, *fq2, *fq2, *fq2) {
	// R ← [2]R
	t0 := new(fq2)
	a := new(fq2).Sqr(&r.x)
	b := new(fq2).Sqr(&r.y)
	c := new(fq2).Sqr(b)
	d := new(fq2).Add(&r.x, b)
	d.Sqr(d).Sub(d, t0.Add(a, c)).Add(d, d)
	e := new(fq2).Add(a, a)
	e.Add(e, a)
	g := new(fq2).Sqr(e)
	sum := new(twistPoint)
	sum.x.Sub(g, t0.Add(d, d))
	sum.y.Sub(d, &sum.x).Mul(&sum.y, e).Sub(&sum.y, t0.Add(t0, t0.Add(t0, t0.Add(c, c))))
	sum.z.Add(&r.y, &r.z).Sqr(&sum.z).Sub(&sum.z, t0.Add(b, &r.t))
	sum.t.Sqr(&sum.z)

	// line function
	c0 := new(fq2).Add(&r.x, e)
	c0.Sqr(c0).Sub(c0, t0.Add(b, b).Add(t0, t0).Add(t0, a).Add(t0, g))
	c1 := new(fq2).Add(e, e)
	c1.Mul(c1, &r.t).Neg(c1)
	fqMul(&c1.c0, &c1.c0, &q.x)
	fqMul(&c1.c1, &c1.c1, &q.x)
	c2 := new(fq2).Mul(&sum.z, &r.t)
	c2.Add(c2, c2)
	fqMul(&c2.c0, &c2.c0, &q.y)
	fqMul(&c2.c1, &c2.c1, &q.y)

	return r.Set(sum), c0, c1, c2
}

// mixedAdditionAndLine returns the sum r + p and the line function result.
// See https://arxiv.org/pdf/0904.0854v3.pdf - Mixed Addition.
func mixedAdditionAndLine(r *twistPoint, p *twistPoint, q *curvePoint, r2 *fq2) (*twistPoint, *fq2, *fq2, *fq2) {
	// R ← R + P
	t0 := new(fq2)
	b := new(fq2).Mul(&p.x, &r.t)
	d := new(fq2).Add(&p.y, &r.z)
	d.Sqr(d).Sub(d, t0.Add(r2, &r.t)).Mul(d, &r.t)
	h := new(fq2).Sub(b, &r.x)
	i := new(fq2).Sqr(h)
	e := new(fq2).Add(i, i)
	e.Add(e, e)
	j := new(fq2).Mul(h, e)
	l1 := new(fq2).Sub(d, t0.Add(&r.y, &r.y))
	v := new(fq2).Mul(&r.x, e)
	sum := new(twistPoint)
	sum.x.Sqr(l1).Sub(&sum.x, t0.Add(v, v).Add(t0, j))
	sum.y.Sub(v, &sum.x).Mul(&sum.y, l1).Sub(&sum.y, t0.Add(&r.y, &r.y).Mul(t0, j))
	sum.z.Add(&r.z, h).Sqr(&sum.z).Sub(&sum.z, t0.Add(&r.t, i))
	sum.t.Sqr(&sum.z)

	// line function
	t1 := new(fq2).Add(l1, l1) // caches 2L1
	c0 := new(fq2).Add(r2, &sum.t)
	c0.Add(c0, t0.Mul(t1, &p.x)).Sub(c0, t0.Add(&p.y, &sum.z).Sqr(t0))
	c1 := new(fq2).Neg(t1)
	fqMul(&c1.c0, &c1.c0, &q.x)
	fqMul(&c1.c1, &c1.c1, &q.x)

	c2 := new(fq2).Add(&sum.z, &sum.z)
	fqMul(&c2.c0, &c2.c0, &q.y)
	fqMul(&c2.c1, &c2.c1, &q.y)

	return r.Set(sum), c0, c1, c2
}

// finalExp implements the final exponentiation step.
// See https://eprint.iacr.org/2019/077.pdf - Algorithm 1, step 8.
func finalExp(p *fq12) *fq12 {
	// easy part
	f := new(fq12).Conjugate(p) // frobenius
	t0 := new(fq12).Inv(p)
	f.Mul(f, t0).Mul(f, t0.Frobenius(f, 2))

	// hard part
	// note: u is negative.
	// See https://eprint.iacr.org/2016/130.pdf - Algorithm 2.
	t0.Sqr(f)
	t1 := new(fq12).Exp(t0, bigU)
	t1.Conjugate(t1)
	t2 := new(fq12).Exp(t1, bigHalfU)
	t2.Conjugate(t2)
	t3 := new(fq12).Conjugate(f)
	t1.Mul(t3, t1).Conjugate(t1).Mul(t1, t2)
	t2.Exp(t1, bigU).Conjugate(t2)
	t3.Exp(t2, bigU).Conjugate(t3)
	t1.Conjugate(t1)
	t3.Mul(t1, t3)
	t1.Conjugate(t1).Frobenius(t1, 3)
	t2.Frobenius(t2, 2)
	t1.Mul(t1, t2)
	t2.Exp(t3, bigU).Conjugate(t2).Mul(t2, t0).Mul(t2, f)
	t1.Mul(t1, t2)
	t2.Frobenius(t3, 1)
	return t1.Mul(t1, t2)
}

// miller implements the Miller’s double-and-add algorithm.
// https://eprint.iacr.org/2016/130.pdf contains useful examples.
func miller(p *curvePoint, q *twistPoint) *fq12 {
	pAffine := new(curvePoint).Set(p).ToAffine()
	qAffine := new(twistPoint).Set(q).ToAffine()
	r := new(twistPoint).Set(qAffine)
	f := new(fq12).SetOne()

	// See https://arxiv.org/pdf/0904.0854v3.pdf - Full addition (precompute R2)
	r2 := new(fq2).Sqr(&qAffine.y)

	// log2(u) - 1 = uArrLen - 2
	for i := uArrLen - 2; i >= 0; i-- {
		// skip initial multiplciation (f = 1)
		if i != (uArrLen - 2) {
			f.Sqr(f)
		}

		_, c0, c1, c4 := doublingAndLine(r, pAffine)
		f.SparseMul014(f, c0, c1, c4)

		if uArr[i] == 1 {
			_, c0, c1, c4 := mixedAdditionAndLine(r, qAffine, pAffine, r2)
			f.SparseMul014(f, c0, c1, c4)
		}
	}

	// u is negative
	return f.Conjugate(f)
}

// Pair implements the optimal ate pairing algorithm on BLS curves.
// See https://eprint.iacr.org/2019/077.pdf - Algorithm 1.
func Pair(g1 *g1Point, g2 *g2Point) *fq12 {
	return finalExp(miller(&g1.p, &g2.p))
}
