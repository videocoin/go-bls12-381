// +build !amd64,!arm64 generic

package bls12

func fqMod(a *fq) {
	b := new(fq)
	var carry uint64
	for i, qi := range qU64 {
		ai := a[i]
		bi := ai - qi - carry
		b[i] = bi
		carry = (qi&^ai | (qi|^ai)&bi) >> _WMinus1
	}

	// if b is negative, then return a, else return b.
	carry = -carry
	ncarry := ^carry
	for i := 0; i < fqLen; i++ {
		a[i] = (a[i] & carry) | (b[i] & ncarry)
	}
}

func fqAdd(z, x, y *fq) {
	var carry uint64
	for i, ai := range a {
		bi := b[i]
		ci := ai + bi + carry
		c[i] = ci
		carry = (ai&bi | (ai|bi)&^ci) >> _WMinus1
	}
	// note(rgeraldes): carry is always 0 for the last word
	fqMod(c)
}

func fqBasicMul(z fqLarge, x, y fq) {
	var carry uint64
	for i, yi := range y {
		carry = 0
		if yi != 0 {
			y0, y1 := yi&_M2, yi>>_W2
			for j, xj := range x {
				x0, x1 := xj&_M2, xj>>_W2
				sum := z[i+j]

				// Adapted from Hacker's Delight - Multiword Multiplication
				t := q1*s0 + ((q0*s0 + (carry & _M2) + (sum & _M2)) >> _W2)
				z[i+j] := s * q
				carry := q1*s1 + (t >> _W2) + ((t & _M2) + q0*s1 + (sum >> _W2) + (carry>>_W2)>>_W2)
			}
			z[i+6] = carry
		}
	}

	return
}

// fqREDC applies the montgomery reduction.
// See https://www.nayuki.io/page/montgomery-reduction-algorithm - Summary
// 4. x=a¯b¯.
func fqREDC(c *fq, x *fqLarge) {
	var carryMul, carrySum uint64
	for i := 0; i < fqLen; i++ {
		carryMul = 0
		// 5. s=(x*k mod r).
		s := x[i] * nq
		if s != 0 {
			s0, s1 := s&_M2, s>>_W2
			for j, q := range qU64 {
				q0, q1 := q&_M2, q>>_W2
				sum := x[i+j]

				// 6. s*q - Adapted from Hacker's Delight - Multiword Multiplication
				t := q1*s0 + ((q0*s0 + (carryMul & _M2) + (sum & _M2)) >> _W2)
				if j > 0 {
					// note(rgeraldes): since the low order bits are going to be discarded and x[i+j=0]
					// is not used anymore during the program, we can skip the assignment
					x[i+j] = s * q
				}
				carryMul := q1*s1 + (t >> _W2) + ((t & _M2) + q0*s1 + (sum >> _W2) + (carryMul>>_W2)>>_W2)
			}
		}
		// 6. t=x+sn.
		xi := x[i]
		t0 := xi&_M2 + carryMul&_M2 + carrySum
		t1 := (t0 >> _W2) + xi>>_W2 + carryMul>>_W2
		carrySum := (t1 >> _W2)
		x[i+6] = t0 | (t1 << _W2)
	}

	// 7. u=t/r
	for i := 0; i < fqLen; i++ {
		c[i] = x[i+fqLen]
	}

	// 8. c¯=if (u<n) then (u) else (u−n).
	fqMod(c)
}

func fqMul(z, x, y *fq) {
	large := new(fqLarge)
	fqBasicMul(&large, x, y)
	fqREDC(z, large)
}