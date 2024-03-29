// +build !amd64,!arm64 generic

package bls12

const (
	halfWordSize = wordSize / 2
	halfWordMask = (1 << halfWordSize) - 1
)

func fqMod(a *fq) {
	b := new(fq)
	var carry uint64
	for i, qi := range q64 {
		ai := a[i]
		bi := ai - qi - carry
		b[i] = bi
		carry = (qi&^ai | (qi|^ai)&bi) >> (wordSize - 1)
	}

	// if b is negative, then return a, else return b.
	carry = -carry
	ncarry := ^carry
	for i := 0; i < fqLen; i++ {
		a[i] = (a[i] & carry) | (b[i] & ncarry)
	}
}

// note(rgeraldes): carry is always 0 for the last word
func fqAdd(z, x, y *fq) {
	var carry uint64
	for i, xi := range x {
		yi := y[i]
		zi := xi + yi + carry
		z[i] = zi
		carry = (xi&yi | (xi|yi)&^zi) >> (wordSize - 1)
	}
	fqMod(z)
}

func fqSub(z, x, y *fq) {
	negY := new(fq)
	fqNeg(negY, y)
	fqAdd(z, x, negY)
}

func fqNeg(z, x *fq) {
	var carry uint64
	for i, qi := range q64 {
		xi := x[i]
		zi := qi - xi - carry
		z[i] = zi
		carry = (xi&^qi | (xi|^qi)&zi) >> (wordSize - 1)
	}
	fqMod(z)
}

func fqBasicMul(z *fqLarge, x, y *fq) {
	var carry uint64
	for i, yi := range y {
		carry = 0
		if yi != 0 {
			y0, y1 := yi&halfWordMask, yi>>halfWordSize
			for j, xj := range x {
				x0, x1 := xj&halfWordMask, xj>>halfWordSize

				// See Hacker's Delight - Multiword Multiplication
				x0y0 := x0 * y0
				x1y0 := x1 * y0
				x0y1 := x0 * y1
				x1y1 := x1 * y1
				w0 := (x0y0 & halfWordMask) + (z[i+j] & halfWordMask) + (carry & halfWordMask)
				w1 := (w0 >> halfWordSize) + (x0y0 >> halfWordSize) + (z[i+j] >> halfWordSize) + (x1y0 & halfWordMask) + (x0y1 & halfWordMask) + (carry >> halfWordSize)
				w2 := (w1 >> halfWordSize) + (x1y0 >> halfWordSize) + (x0y1 >> halfWordSize) + (x1y1 & halfWordMask)
				carry = (((w2 >> halfWordSize) + (x1y1 >> halfWordSize)) << halfWordSize) | (w2 & halfWordMask)
				z[i+j] = (w1 << halfWordSize) | (w0 & halfWordMask)
			}
			z[i+fqLen] = carry
		}
	}
}

// fqREDC applies the montgomery reduction.
// See https://www.nayuki.io/page/montgomery-reduction-algorithm - Summary
// 4. x=a¯b¯.
func fqREDC(c *fq, x *fqLarge) {
	var carryMul, carrySum uint64
	for i := 0; i < fqLen; i++ {
		carryMul = 0
		// 2. k=(r(r^−1 mod n)−1)/n
		// 5. s=(x*k mod r);
		s := x[i] * qK64
		if s != 0 {
			s0, s1 := s&halfWordMask, s>>halfWordSize
			for j, q := range q64 {
				q0, q1 := q&halfWordMask, q>>halfWordSize

				// See Hacker's Delight - Multiword Multiplication
				q0s0 := q0 * s0
				q1s0 := q1 * s0
				q0s1 := q0 * s1
				q1s1 := q1 * s1
				w0 := (q0s0 & halfWordMask) + (x[i+j] & halfWordMask) + (carryMul & halfWordMask)
				w1 := (w0 >> halfWordSize) + (q0s0 >> halfWordSize) + (x[i+j] >> halfWordSize) + (q1s0 & halfWordMask) + (q0s1 & halfWordMask) + (carryMul >> halfWordSize)
				w2 := (w1 >> halfWordSize) + (q1s0 >> halfWordSize) + (q0s1 >> halfWordSize) + (q1s1 & halfWordMask)
				carryMul = (((w2 >> halfWordSize) + (q1s1 >> halfWordSize)) << halfWordSize) | (w2 & halfWordMask)
				if j > 0 {
					// note(rgeraldes): since the low order bits are going to be discarded and
					// x[i+j=0] is not used anymore during the program, we can skip the assignment.
					x[i+j] = (w1 << halfWordSize) | (w0 & halfWordMask)
				}
			}
		}

		// 6. t=x+sn.
		xi := x[i+fqLen]
		t0 := xi&halfWordMask + carryMul&halfWordMask + carrySum&halfWordMask
		t1 := (t0 >> halfWordSize) + (xi >> halfWordSize) + (carryMul >> halfWordSize) + (carrySum >> halfWordSize)
		carrySum = t1 >> halfWordSize
		x[i+fqLen] = (t0 & halfWordMask) | (t1 << halfWordSize)
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
	fqBasicMul(large, x, y)
	fqREDC(z, large)
}
