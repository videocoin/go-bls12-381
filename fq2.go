package bls12

// fq2 is an element of Fq² = Fq[X]/(X² − β), where β is a quadratic non-residue
// in Fq with a value of -1. See http://eprint.iacr.org/2006/471.pdf for
// arithmetic.
type fq2 struct {
	c0, c1 fq
}

// Set sets z to x and returns z.
func (z *fq2) Set(x *fq2) *fq2 {
	z.c0.Set(&x.c0)
	z.c1.Set(&x.c1)
	return z
}

// SetZero sets z to 0 and returns z.
func (z *fq2) SetZero() *fq2 {
	*z = fq2{}
	return z
}

// SetOne sets z to 1 and returns z.
func (z *fq2) SetOne() *fq2 {
	z.c0 = *new(fq).SetUint64(1)
	z.c1 = fq{}
	return z
}

// Add sets z to the sum x+y and returns z.
func (z *fq2) Add(x, y *fq2) *fq2 {
	fqAdd(&z.c0, &x.c0, &y.c0)
	fqAdd(&z.c1, &x.c1, &y.c1)
	return z
}

// Neg sets z to -x and returns z.
func (z *fq2) Neg(x *fq2) *fq2 {
	fqNeg(&z.c0, &x.c0)
	fqNeg(&z.c1, &x.c1)
	return z
}

// Sub sets z to the difference x-y and returns z.
func (z *fq2) Sub(x, y *fq2) *fq2 {
	fqSub(&z.c0, &x.c0, &y.c0)
	fqSub(&z.c1, &x.c1, &y.c1)
	return z
}

// Mul sets z to the product x*y and returns z. Mul utilizes Karatsuba's method.
func (z *fq2) Mul(x, y *fq2) *fq2 {
	ret := new(fq2)
	// v0 = a0b0
	// v1 = a1b1
	// c0 = v0 + βv1
	// c0 = v0 - v1
	v0, v1 := new(fq), new(fq)
	fqMul(v0, &x.c0, &y.c0)
	fqMul(v1, &x.c1, &y.c1)
	fqSub(&ret.c0, v0, v1)

	// c1 = (a0 + a1)(b0 + b1) − v0 − v1
	// c1 = (a0 + a1)(b0 + b1) − (v0 + v1)
	a, b, v := new(fq), new(fq), new(fq)
	fqAdd(a, &x.c0, &x.c1)
	fqAdd(b, &y.c0, &y.c1)
	fqMul(&ret.c1, a, b)
	fqAdd(v, v0, v1)
	fqSub(&ret.c1, &ret.c1, v)

	return z.Set(ret)
}

// MulXi sets z to the product ξX and returns z.
func (z *fq2) MulXi(x *fq2) *fq2 {
	// ξ = u + 1
	// X = x + yu
	// ξX = (1 + u)(x + yu)
	// ξX = x + yu + xu + yu^2
	// ξX = x + yu + xu - y
	// ξX = x - y + u(x + y)
	ret := new(fq2)
	fqSub(&ret.c0, &x.c0, &x.c1)
	fqAdd(&ret.c1, &x.c0, &x.c1)

	return z.Set(ret)
}

// Sqr sets z to the product x*x and returns z.
// Sqr utilizes complex squaring.
func (z *fq2) Sqr(x *fq2) *fq2 {
	// v0 = a0a1
	v0 := new(fq)
	fqMul(v0, &x.c0, &x.c1)
	// c0 = (a0 + a1)(a0 - a1)
	ret := new(fq2)
	fqAdd(&ret.c0, &x.c0, &x.c1)
	fqSub(&ret.c1, &x.c0, &x.c1)
	fqMul(&ret.c0, &ret.c0, &ret.c1)
	// c1 = 2v0
	fqAdd(&ret.c1, v0, v0)

	return z.Set(ret)
}

// Inv sets z to 1/x and returns z.
func (z *fq2) Inv(x *fq2) *fq2 {
	// t0 = x0^2
	// t1 = x1^2
	// t0 = t0 + t1
	// t0 = 1/t0
	// z0 = x0 * t0
	// z1 = - x1 * t0
	t0, t1 := new(fq), new(fq)
	fqMul(t0, &x.c0, &x.c0)
	fqMul(t1, &x.c1, &x.c1)
	fqAdd(t0, t0, t1)
	fqInv(t0, t0)
	fqMul(&z.c0, &x.c0, t0)
	fqMul(&z.c1, &x.c1, t0)
	fqNeg(&z.c1, &z.c1)

	return z
}

func (z *fq2) Frobenius(x *fq2, power uint64) *fq2 {
	fqMul(&z.c1, &x.c1, frobFq2C1[power%2])
	return z
}
