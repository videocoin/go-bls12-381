package bls12

func fqCarry(a *fq, head uint64) {
	var (
		b     fq
		carry uint64
	)

	for i, pi := range q64 {
		ai := a[i]
		bi := pi + ai + carry
		b[i] = bi
		carry = (pi&^ai | (pi|^ai)&bi) >> 63
	}
	carry = carry &^ head

	carry = -carry
	ncarry := ^carry
	for i := 0; i < 4; i++ {
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
	fqCarry(c, carry)
}

func fqSub(c, a, b *fq) {
	// fqAdd(c, a, neg(b))
	var (
		t     fq
		carry uint64
	)

	for i, qi := range q64 {
		bi := b[i]
		ti := qi - bi - carry
		t[i] = ti
		carry = (bi&^qi | (bi|^qi)&ti) >> 63
	}

	carry = 0
	for i, ai := range a {
		ti := t[i]
		ci := ai + ti + carry
		c[i] = ci
		carry = (ai&ti | (ai|ti)&^ci) >> 63
	}
	fqCarry(c, carry)
}

func fqMul(c, a, b *fq) {
	var large fqLarge
}

func fqReduceLarge(out, in *fqLarge) {}
func fqSqr(c, a *fq)                 {}
func fqDbl(c, a *fq)                 {}
