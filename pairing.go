package bls12

// doublingAndLine sets z to the sum z + t and f to the line function result and returns
// the pair (z, f). See https://arxiv.org/pdf/0904.0854v3.pdf - Doubling on curves
// with a4 = 0. TODO: q must be affine?
func doublingAndLine(r *twistPoint, q *curvePoint) (*twistPoint, *fq12) {
	// R ← [2]R
	// note: there's a faster way to compute the doubling (2m + 5s instead of 1m +
	// 7s) but line functions make use of T1 = Z². TODO: benchmark variable allocation
	ret := new(twistPoint)
	a := new(fq2).Sqr(&r.x)
	b := new(fq2).Sqr(&r.y)
	c := new(fq2).Sqr(b)
	d := new(fq2).Dbl(new(fq2).Sub(new(fq2).Sqr(new(fq2).Add(&r.x, b)), new(fq2).Add(a, c)))
	e := new(fq2).Add(new(fq2).Dbl(a), a)
	g := new(fq2).Sqr(e)
	ret.x.Sub(g, new(fq2).Dbl(d))
	ret.y.Sub(new(fq2).Mul(e, new(fq2).Sub(d, &ret.x)), new(fq2).Dbl(new(fq2).Dbl(new(fq2).Dbl(c))))
	ret.z.Sub(new(fq2).Sqr(new(fq2).Add(&r.y, &r.z)), new(fq2).Add(b, &r.t))
	ret.t.Sqr(&ret.z)

	// line function
	/*
		l := new(fq12)
	*/

	return ret, &fq12{}
}

// mixedAdditionAndLine sets z to the sum z + t and f to the line function result and returns
// the pair (z, f). See https://arxiv.org/pdf/0904.0854v3.pdf - Mixed Addition.
// Mixed addition means that the second input point is in affine representation.
func mixedAdditionAndLine(r *twistPoint, p *twistPoint, q *curvePoint, r2 *fq2) (*twistPoint, *fq12) {
	// R ← R + P
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
	ret := new(twistPoint)

	// TODO
	ret.x = *new(fq2).Sub(new(fq2).Sqr(l1), new(fq2).Add(j, new(fq2).Dbl(v)))
	// y3 = r · (V − X3) − 2Y1 · J;
	//ret.y =
	// Z3 = (Z1 + H)^2 − T1 − I
	ret.z.Add(&r.z, h).Sqr(&ret.z).Sub(&ret.z, new(fq2).Add(&r.t, i))
	ret.t.Sqr(&ret.z)

	// line function
	/*
		l := new(fq12)
		a := new(fq2).Dbl(&ret.z)
		//a.Mul(a, q.)
		b := new(fq2).Add(r2, &ret.t)
	*/

	return ret, &fq12{}
}

/*
http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.215.7255&rep=rep1&type=pdf

The first two parts of the exponentiation are “easy” as raising to the power of
p is an almost free application of the Frobenius operator, as p is the characteristic
of the extension field. However the first part of the exponentiation is not only
cheap (although it does require an extension field division), it also simplifies the
rest of the final exponentiation. After raising to the power of p
d − 1 the field element becomes “unitary” [20]. This has important implications, as squaring of
unitary elements is significantly cheaper than squaring of non-unitary elements,
and any future inversions can be implemented by simple conjugation [21], [20],
[12].

*/
// k = 12; (p^12 − 1)/r = (p^6 − 1).(p^2 + 1).[(p^4 − p^2 + 1)/r]
// p^6-power Frobenius automorphism on Fp12 ,a single inversion and a multiplication in Fp12
func finalExp(p *fq12) *fq12 {
	// raising to the power of p - frobenius operator
	// It can be done by applying the p^6 -power Frobenius automorphism on fq12
	// The automorphism maps every element to its p-th power
	return &fq12{}
}

// miller implements the Miller’s double-and-add algorithm. Non Adjacent Form
// does not reduce the number of additions for this specific value of u.
func miller(p *g1Point, q *g2Point) *fq12 {
	f := new(fq12).SetOne()
	//t := new(g2Point).Set(q).ToAffine()

	for i := log2U; i < 0; i++ {
		f.Sqr(f)
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
