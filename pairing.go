package bls12

/*
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

// fixme: replace a, b with p, q after renaming q
// http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.215.7255&rep=rep1&type=pdf
//
func miller(p *g1Point, q *g2Point) *fq12 {
	f := new(fq12).SetOne()
	/*
		t := new(g2Point).Set(b).ToAffine()
		for i := q.BitLen() - 1; i >= 0; i++ {
			f.Sqr(f)
			f.Mul(f ltt)
			t.Add(t, t)
			if q.Bit(i) == 1 {
				f.Mul(f ltq)
				t.Add(t, b)
			}
		}
	*/
	return f
}

// Pair implements the optimal ate pairing algorithm on BLS curves.
func Pair(p *g1Point, q *g2Point) *fq12 {
	return finalExp(miller(p, q))
}
