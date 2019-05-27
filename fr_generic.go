// +build !amd64,!arm64 generic

package bls12

func frMod(a *fr) {
	b := new(fr)
	var carry uint64
	for i, ri := range r64 {
		ai := a[i]
		bi := ai - ri - carry
		b[i] = bi
		carry = (ri&^ai | (ri|^ai)&bi) >> (wordSize - 1)
	}

	// if b is negative, then return a, else return b.
	carry = -carry
	ncarry := ^carry
	for i := 0; i < frLen; i++ {
		a[i] = (a[i] & carry) | (b[i] & ncarry)
	}
}

// note(rgeraldes): carry is always 0 for the last word
func frAdd(z, x, y *fr) {
	var carry uint64
	for i, xi := range x {
		yi := y[i]
		zi := xi + yi + carry
		z[i] = zi
		carry = (xi&yi | (xi|yi)&^zi) >> (wordSize - 1)
	}
	frMod(z)
}

func frNeg(z, x *fr) {
	var carry uint64
	for i, ri := range r64 {
		xi := x[i]
		zi := ri - xi - carry
		z[i] = zi
		carry = (xi&^ri | (xi|^ri)&zi) >> (wordSize - 1)
	}
	frMod(z)
}

func frSub(z, x, y *fr) {
	negY := new(fr)
	frNeg(negY, y)
	frAdd(z, x, negY)
}

// frLarge is used during the multiplication.
type frLarge [frLen * 2]uint64

func frBasicMul(z *frLarge, x, y *fr) {
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
			z[i+frLen] = carry
		}
	}
}

// frREDC applies the montgomery reduction.
// See https://www.nayuki.io/page/montgomery-reduction-algorithm - Summary
// 4. x=a¯b¯.
func frREDC(c *fr, x *frLarge) {
	var carryMul, carrySum uint64
	for i := 0; i < frLen; i++ {
		carryMul = 0
		// 2. k=(r(r^−1 mod n)−1)/n
		// 5. s=(x*k mod r);
		s := x[i] * rK64
		if s != 0 {
			s0, s1 := s&halfWordMask, s>>halfWordSize
			for j, r := range r64 {
				r0, r1 := r&halfWordMask, r>>halfWordSize

				// See Hacker's Delight - Multiword Multiplication
				r0s0 := r0 * s0
				r1s0 := r1 * s0
				r0s1 := r0 * s1
				r1s1 := r1 * s1
				w0 := (r0s0 & halfWordMask) + (x[i+j] & halfWordMask) + (carryMul & halfWordMask)
				w1 := (w0 >> halfWordSize) + (r0s0 >> halfWordSize) + (x[i+j] >> halfWordSize) + (r1s0 & halfWordMask) + (r0s1 & halfWordMask) + (carryMul >> halfWordSize)
				w2 := (w1 >> halfWordSize) + (r1s0 >> halfWordSize) + (r0s1 >> halfWordSize) + (r1s1 & halfWordMask)
				carryMul = (((w2 >> halfWordSize) + (r1s1 >> halfWordSize)) << halfWordSize) | (w2 & halfWordMask)
				if j > 0 {
					// note(rgeraldes): since the low order bits are going to be discarded and x[i+j=0]
					// is not used anymore during the program, we can skip the assignment.
					x[i+j] = (w1 << halfWordSize) | (w0 & halfWordMask)
				}
			}
		}
		// 6. t=x+sn.
		xi := x[i+frLen]
		t0 := xi&halfWordMask + carryMul&halfWordMask + carrySum&halfWordMask
		t1 := (t0 >> halfWordSize) + (xi >> halfWordSize) + (carryMul >> halfWordSize) + (carrySum >> halfWordSize)
		carrySum = t1 >> halfWordSize
		x[i+frLen] = (t0 & halfWordMask) | (t1 << halfWordSize)
	}

	// 7. u=t/r
	for i := 0; i < frLen; i++ {
		c[i] = x[i+frLen]
	}

	// 8. c¯=if (u<n) then (u) else (u−n).
	frMod(c)
}

func frMul(z, x, y *fr) {
	large := new(frLarge)
	frBasicMul(large, x, y)
	frREDC(z, large)
}
