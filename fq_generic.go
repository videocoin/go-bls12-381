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

func fqAdd(c, a, b *fq) {
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
