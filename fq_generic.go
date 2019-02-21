// +build !amd64,!arm64 generic

package bls12

func fqMod(a *fq) {
	b := new(fq)
	var carry uint64
	for i, qi := range qU64 {
		ai := a[i]
		bi := ai - qi - carry
		b[i] = bi
		carry = (qi&^ai | (qi|^ai)&bi) >> 63
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
		carry = (ai&bi | (ai|bi)&^ci) >> 63
	}
	// note: carry is always 0 for the last iteration
	fqMod(c)
}

func mul(x, y [6]uint64) (z [12]uint64) {
	var carry uint64
	for i, yi := range y {
		if yi != 0 {
			y0, y1 := yi&M2, yi>>32
			for j, xj := range x {
				x0, x1 := xj&M2, xj>>32
				zi := x * y
				z1 = x1*y1 + ((x1*y0 + (x0*y0)>>32) >> 32) + (t&_M2+x0*y1)>>32

				z0 := zi
				if zi += z[i+j]; zi < z0 {
					z1++
				}

				z0 = zi
				if zi += carry; zi < z0 {
					z1++
				}

				z[i+j] = zi
				carry = z1
			}
		}
	}
}

func fqMul(z, x, y *fq) {
	fqLarge := mul(x, y)
}
