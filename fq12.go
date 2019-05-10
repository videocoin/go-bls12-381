package bls12

import "math/big"

var (
	// Fq2(u + 1)**(((p^power) - 1) / 6), power E [0, 11]
	frob12c1 = [12]*fq2{
		fq2One,
		&fq2{
			fq{0x7089552b319d465, 0xc6695f92b50a8313, 0x97e83cccd117228f, 0xa35baecab2dc29ee, 0x1ce393ea5daace4d, 0x8f2220fb0fb66eb},
			fq{0xb2f66aad4ce5d646, 0x5842a06bfc497cec, 0xcf4895d42599d394, 0xc11b9cba40a8e8d0, 0x2e3813cbe5a0de89, 0x110eefda88847faf},
		},
		&fq2{
			fq{0xecfb361b798dba3a, 0xc100ddb891865a2c, 0xec08ff1232bda8e, 0xd5c13cc6f1ca4721, 0x47222a47bf7b5c04, 0x110f184e51c5f59},
			fq{},
		},
		&fq2{
			fq{0x3e2f585da55c9ad1, 0x4294213d86c18183, 0x382844c88b623732, 0x92ad2afd19103e18, 0x1d794e4fac7cf0b9, 0xbd592fc7d825ec8},
			fq{0x7bcfa7a25aa30fda, 0xdc17dec12a927e7c, 0x2f088dd86b4ebef1, 0xd1ca2087da74d4a7, 0x2da2596696cebc1d, 0xe2b7eedbbfd87d2},
		},
		&fq2{
			fq{0x30f1361b798a64e8, 0xf3b8ddab7ece5a2a, 0x16a8ca3ac61577f7, 0xc26a2ff874fd029b, 0x3636b76660701c6e, 0x51ba4ab241b6160},
			fq{},
		},
		&fq2{
			fq{0x3726c30af242c66c, 0x7c2ac1aad1b6fe70, 0xa04007fbba4b14a2, 0xef517c3266341429, 0x95ba654ed2226b, 0x2e370eccc86f7dd},
			fq{0x82d83cf50dbce43f, 0xa2813e53df9d018f, 0xc6f0caa53c65e181, 0x7525cf528d50fe95, 0x4a85ed50f4798a6b, 0x171da0fd6cf8eebd},
		},
		&fq2{
			fq{0x43f5fffffffcaaae, 0x32b7fff2ed47fffd, 0x7e83a49a2e99d69, 0xeca8f3318332bb7a, 0xef148d1ea0f4c069, 0x40ab3263eff0206},
			fq{},
		},
		&fq2{
			fq{0xb2f66aad4ce5d646, 0x5842a06bfc497cec, 0xcf4895d42599d394, 0xc11b9cba40a8e8d0, 0x2e3813cbe5a0de89, 0x110eefda88847faf},
			fq{0x7089552b319d465, 0xc6695f92b50a8313, 0x97e83cccd117228f, 0xa35baecab2dc29ee, 0x1ce393ea5daace4d, 0x8f2220fb0fb66eb},
		},
		&fq2{
			fq{0xcd03c9e48671f071, 0x5dab22461fcda5d2, 0x587042afd3851b95, 0x8eb60ebe01bacb9e, 0x3f97d6e83d050d2, 0x18f0206554638741},
			fq{},
		},
		&fq2{
			fq{0x7bcfa7a25aa30fda, 0xdc17dec12a927e7c, 0x2f088dd86b4ebef1, 0xd1ca2087da74d4a7, 0x2da2596696cebc1d, 0xe2b7eedbbfd87d2},
			fq{0x3e2f585da55c9ad1, 0x4294213d86c18183, 0x382844c88b623732, 0x92ad2afd19103e18, 0x1d794e4fac7cf0b9, 0xbd592fc7d825ec8},
		},
		&fq2{
			fq{0x890dc9e4867545c3, 0x2af322533285a5d5, 0x50880866309b7e2c, 0xa20d1b8c7e881024, 0x14e4f04fe2db9068, 0x14e56d3f1564853a},
			fq{},
		},
		&fq2{
			fq{0x82d83cf50dbce43f, 0xa2813e53df9d018f, 0xc6f0caa53c65e181, 0x7525cf528d50fe95, 0x4a85ed50f4798a6b, 0x171da0fd6cf8eebd},
			fq{0x3726c30af242c66c, 0x7c2ac1aad1b6fe70, 0xa04007fbba4b14a2, 0xef517c3266341429, 0x95ba654ed2226b, 0x2e370eccc86f7dd},
		},
	}
)

// fq12 implements the field of size q¹² as a quadratic extension of fq6 where
// y = v. See https://eprint.iacr.org/2006/471.pdf for arithmetic.
type fq12 struct {
	c0, c1 fq6
}

// Set sets z to x and returns z.
func (z *fq12) Set(x *fq12) *fq12 {
	z.c0.Set(&x.c0)
	z.c1.Set(&x.c1)
	return z
}

// SetOne sets z to 1 and returns z.
func (z *fq12) SetOne() *fq12 {
	z.c0.SetOne()
	z.c1.SetZero()
	return z
}

// Equal reports whether x is equal to y.
func (x *fq12) Equal(y *fq12) bool {
	return *x == *y
}

// Conjugate sets z to the conjugate of x and returns z.
func (z *fq12) Conjugate(x *fq12) *fq12 {
	z.c0.Set(&x.c0)
	z.c1.Neg(&x.c1)
	return z
}

// Add sets z to the sum x+y and returns z.
func (z *fq12) Add(x, y *fq12) *fq12 {
	z.c0.Add(&x.c0, &y.c0)
	z.c1.Add(&x.c1, &y.c1)

	return z
}

// Mul sets z to the product x*y and returns z.
// Mul utilizes Karatsuba's method.
func (z *fq12) Mul(x, y *fq12) *fq12 {
	ret, t0, v1 := new(fq12), new(fq6), new(fq6).Mul(&x.c1, &y.c1)
	ret.c0.MulQuadraticNonResidue(v1).Add(&ret.c0, v1)
	ret.c1.Add(&x.c0, &x.c1).Mul(&ret.c1, t0.Add(&y.c0, &y.c1)).Sub(&ret.c1, t0.Add(t0.Mul(&x.c0, &y.c0), v1))

	return z.Set(ret)
}

// SparseMult sets z to the product of x with a0, a1, a2 and returns z.
// SparseMult utilizes the sparness property to avoid full fq12 arithmetic.
// See https://eprint.iacr.org/2012/408.pdf - algorithm 5.
func (z *fq12) SparseMult(x *fq12, a0 *fq2, a1 *fq2, a2 *fq2) *fq12 {
	a, t0 := new(fq6), new(fq6)
	a.c0.Mul(a0, &x.c0.c0)
	a.c1.Mul(a0, &x.c0.c1)
	a.c2.Mul(a0, &x.c0.c2)
	b := new(fq6).SparseMul(&x.c1, a1, a2)
	c := new(fq6)
	c.c0.Add(a0, a1)
	c.c1 = *a2
	e := new(fq6).SparseMul(t0.Add(&x.c0, &x.c1), &c.c0, &c.c1) // d, e
	z.c1.Sub(e, t0.Add(a, b))                                   // f
	z.c0.Add(a, t0.MulQuadraticNonResidue(b))                   // g, h

	return z
}

// Sqr sets z to the product x*x and returns z.
// Sqr utilizes complex squaring.
func (z *fq12) Sqr(x *fq12) *fq12 {
	ret, t0, v0 := new(fq12), new(fq6), new(fq6).Mul(&x.c0, &x.c1)
	ret.c0.Add(&x.c0, &x.c1).Mul(&ret.c0, t0.MulQuadraticNonResidue(&x.c1).Add(t0, &x.c0)).Sub(&ret.c0, t0.Add(v0, t0.MulQuadraticNonResidue(v0)))
	ret.c1.Add(v0, v0)

	return z.Set(ret)
}

// Inv sets z to 1/x and returns z.
// See "Implementing cryptographic pairings", M. Scott - section 3.2.
func (z *fq12) Inv(x *fq12) *fq12 {
	t0, t1 := new(fq6).Sqr(&x.c0), new(fq6).Sqr(&x.c1)
	t1.MulQuadraticNonResidue(t1).Sub(t0, t1).Inv(t1)
	z.c0.Mul(&x.c0, t1)
	z.c1.Mul(&x.c1, t1)

	return z
}

// Exp sets z=x**y and returns z.
func (z *fq12) Exp(x *fq12, y *big.Int) *fq12 {
	ret, base := new(fq12).SetOne(), *x
	for i := y.BitLen() - 1; i >= 0; i-- {
		if y.Bit(i) == 1 {
			ret.Mul(ret, &base)
		}
		base.Sqr(&base)
	}

	return z.Set(ret)
}

// Frobenius sets z to the pth-power Frobenius of x and returns z.
func (z *fq12) Frobenius(x *fq12, power uint64) *fq12 {
	ret := new(fq12)
	ret.c0.Frobenius(&x.c0, power)
	ret.c1.Frobenius(&x.c1, power)
	ret.c1.c0.Mul(&ret.c1.c0, frob12c1[power%12])
	ret.c1.c1.Mul(&ret.c1.c1, frob12c1[power%12])
	ret.c1.c2.Mul(&ret.c1.c2, frob12c1[power%12])

	return z.Set(ret)
}
