// +build !amd64,!arm64 generic

package bls12

import "fmt"

const (
	wordSize     = 64
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

func fqAdd(z, x, y *fq) {
	var carry uint64
	for i, xi := range x {
		yi := y[i]
		zi := xi + yi + carry
		z[i] = zi
		carry = (xi&yi | (xi|yi)&^zi) >> (wordSize - 1)
	}
	// note(rgeraldes): carry is always 0 for the last word
	fqMod(z)
}

func fqDbl(z, x *fq) {
	fqAdd(z, x, x)
}

func fqSub(z, x, y *fq) {
	fqNeg(y, y)
	fqAdd(z, x, y)
}

func fqNeg(z, x *fq) {
	var carry uint64
	for i, qi := range q64 {
		xi := x[i]
		zi := qi - xi - carry
		z[i] = zi
		carry = (xi&^qi | (xi|^qi)&zi) >> 63
	}
}

func fqBasicMul(z *fqLarge, x, y *fq) {
	var carry uint64
	for i, yi := range y {
		carry = 0
		if yi != 0 {
			y0, y1 := yi&halfWordMask, yi>>halfWordSize
			for j, xj := range x {
				x0, x1 := xj&halfWordMask, xj>>halfWordSize
				sum := z[i+j]

				// Adapted from Hacker's Delight - Multiword Multiplication
				t := x1*y0 + ((x0*y0 + (carry & halfWordMask) + (sum & halfWordMask)) >> halfWordSize)
				z[i+j] = xj * yi
				carry = x1*y1 + (t >> halfWordSize) + ((t & halfWordMask) + x0*y1 + (sum >> halfWordSize) + (carry>>halfWordSize)>>halfWordSize)
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
		// 2. k=(r(r^−1 mod n)−1)/n
		// 5. s=(x*k mod r);
		s := x[i] * k64
		if s != 0 {
			s0, s1 := s&halfWordMask, s>>halfWordSize
			for j, q := range q64 {
				q0, q1 := q&halfWordMask, q>>halfWordSize
				sum := x[i+j]

				fmt.Println(s)
				fmt.Println(q)
				fmt.Println(sum)

				// 6. s*q - Adapted from Hacker's Delight - Multiword Multiplication
				t := q1*s0 + ((q0*s0 + (carryMul & halfWordMask) + (sum & halfWordMask)) >> halfWordSize)
				if j > 0 {
					// note(rgeraldes): since the low order bits are going to be discarded and x[i+j=0]
					// is not used anymore during the program, we can skip the assignment.
					x[i+j] = s * q
					fmt.Println(s)
					fmt.Println(q)
					fmt.Println(x[i+j])
				}
				carryMul = q1*s1 + (t >> halfWordSize) + ((t & halfWordMask) + q0*s1 + (sum >> halfWordSize) + (carryMul>>halfWordSize)>>halfWordSize)
				fmt.Println(carryMul)
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

func fqSqr(z, x *fq) {
	fqMul(z, x, x)
}
