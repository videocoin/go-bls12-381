package bls12

func fqCarry(a *fq, head uint64) {
	var b fq
	var carry uint64
	for i, pi := range q {
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

func fqMul(c, a, b *fq) {

}
