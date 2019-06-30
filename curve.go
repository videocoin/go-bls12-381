package bls12

import (
	"math/big"

	"golang.org/x/crypto/blake2b"
)

const (
	compressedFormMask  uint8 = 1 << 7
	pointAtInfinityMask uint8 = 1 << 6
)

var (
	fqCurveB, _                   = new(fq).SetString("4")
	fqCurveBPlusOne, _            = new(fq).SetString("5")
	fqSqrtNegThree, _             = new(fq).SetString("1586958781458431025242759403266842894121773480562120986020912974854563298150952611241517463240701")
	fqHalfSqrtNegThreeMinusOne, _ = new(fq).SetString("793479390729215512621379701633421447060886740281060493010456487427281649075476305620758731620350")

	iso11XNum = []*fq{}
	iso11XDen = []*fq{}
	iso11YNum = []*fq{}
	iso11YDen = []*fq{}

	// Values taken from the execution of https://eprint.iacr.org/2019/403.pdf - A The isogeny maps.
	iso11K = [][]*fq{iso11XNum, iso11XDen, iso11YNum, iso11YDen}
)

// curvePoint is an elliptic curve point in projective coordinates. The elliptic
// curve is defined by the following equation y²=x³+3.
type curvePoint struct {
	x, y, z fq
}

// Set sets c to the value of a and returns c.
func (c *curvePoint) Set(a *curvePoint) *curvePoint {
	c.x, c.y, c.z = a.x, a.y, a.z
	return c
}

// Equal reports whether a is equal to b.
func (a *curvePoint) Equal(b *curvePoint) bool {
	return *a == *b
}

// IsInfinity reports whether the point is at infinity.
func (a *curvePoint) IsInfinity() bool {
	return a.z == fq{}
}

// Add sets c to the sum a+b and returns c.
func (c *curvePoint) Add(a, b *curvePoint) *curvePoint {
	if a.IsInfinity() {
		return c.Set(b)
	}
	if b.IsInfinity() {
		return c.Set(a)
	}

	// faster than Add
	if a.Equal(b) {
		return c.Double(a)
	}

	// See https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	z1z1, z2z2 := new(fq), new(fq)
	fqMul(z1z1, &a.z, &a.z)
	fqMul(z2z2, &b.z, &b.z)

	u1, u2 := new(fq), new(fq)
	fqMul(u1, &a.x, z2z2)
	fqMul(u2, &b.x, z1z1)

	s1, s2 := new(fq), new(fq)
	fqMul(s1, &a.y, &b.z)
	fqMul(s1, s1, z2z2)
	fqMul(s2, &b.y, &a.z)
	fqMul(s2, s2, z1z1)

	h, i, j, r, v := new(fq), new(fq), new(fq), new(fq), new(fq)
	fqSub(h, u2, u1)
	fqAdd(i, h, h)
	fqMul(i, i, i)
	fqMul(j, h, i)
	fqSub(r, s2, s1)
	fqAdd(r, r, r)
	fqMul(v, u1, i)

	p, t0, t1 := new(curvePoint), new(fq), new(fq)
	fqAdd(t0, v, v)
	fqAdd(t0, t0, j)
	fqMul(&p.x, r, r)
	fqSub(&p.x, &p.x, t0)

	fqAdd(t0, s1, s1)
	fqMul(t0, t0, j)
	fqSub(t1, v, &p.x)
	fqMul(t1, t1, r)
	fqSub(&p.y, t1, t0)

	fqAdd(&p.z, &a.z, &b.z)
	fqMul(&p.z, &p.z, &p.z)
	fqAdd(t0, z1z1, z2z2)
	fqSub(&p.z, &p.z, t0)
	fqMul(&p.z, &p.z, h)

	return c.Set(p)
}

// Double sets c to the sum a+a and returns c.
// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
func (c *curvePoint) Double(a *curvePoint) *curvePoint {
	d, e, f, g, h, i := new(fq), new(fq), new(fq), new(fq), new(fq), new(fq)
	fqMul(d, &a.x, &a.x)
	fqMul(e, &a.y, &a.y)
	fqMul(f, e, e)
	fqMul(g, &a.x, e)
	fqAdd(g, g, g)
	fqAdd(g, g, g)
	fqAdd(h, d, d)
	fqAdd(h, h, d)
	fqMul(i, h, h)

	p, t0 := new(curvePoint), new(fq)
	fqAdd(&p.x, g, g)
	fqSub(&p.x, i, &p.x)
	fqAdd(t0, f, f)
	fqAdd(t0, t0, t0)
	fqAdd(t0, t0, t0)
	fqSub(&p.y, g, &p.x)
	fqMul(&p.y, &p.y, h)
	fqSub(&p.y, &p.y, t0)
	fqMul(&p.z, &a.y, &a.z)
	fqAdd(&p.z, &p.z, &p.z)

	return c.Set(p)
}

// ScalarMult returns b*(Ax,Ay) where b is a number in big-endian form.
// See https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add.
func (c *curvePoint) ScalarMult(a *curvePoint, b *big.Int) *curvePoint {
	p := new(curvePoint)
	for i := b.BitLen() - 1; i >= 0; i-- {
		p.Double(p)
		if b.Bit(i) == 1 {
			p.Add(p, a)
		}
	}

	return c.Set(p)
}

// ToAffine sets a to its affine value and returns a.
func (a *curvePoint) ToAffine() *curvePoint {
	if a.z == *new(fq).SetUint64(1) {
		return a
	}

	// TODO create new curve point
	if a.IsInfinity() {
		//  If this bit is set, the remaining bits of the group element's encoding should be set to zero.
		//pointAtInfinityMask
		return nil
	}

	zInv, zInvSqr, zInvCube := new(fq), new(fq), new(fq)
	fqInv(zInv, &a.z)
	fqMul(zInvSqr, zInv, zInv)
	fqMul(zInvCube, zInvSqr, zInv)
	fqMul(&a.x, &a.x, zInvSqr)
	fqMul(&a.y, &a.y, zInvCube)
	a.z = *new(fq).SetUint64(1)

	return a
}

// CompressedEncode converts a curve point into the uncompressed form specified in
// See https://github.com/zkcrypto/pairing/tree/master/src/bls12_381#serialization.
func (a *curvePoint) Marshal() []byte {
	a.ToAffine()

	x := new(fq).MontgomeryDecode(&a.x)
	y := new(fq).MontgomeryDecode(&a.y)

	ret := make([]byte, fqByteLen*2)
	copy(ret, x.Bytes())
	copy(ret[fqByteLen:], y.Bytes())

	// TODO review
	//if cp.IsInfinity() {
	//	ret[0] |= pointAtInfinityMask
	//}

	return ret
}

// Unmarshal decodes a curve point, serialized by Marshal.
// It is an error if the point is not on the curve.
func (cp *curvePoint) Unmarshal(data []byte) error {
	if len(data) != 2*fqByteLen {
		// TODO error
		return nil
	}

	if data[0]&compressedFormMask != 0 { // uncompressed form
		// TODO error
		return nil
	}

	if data[0]&pointAtInfinityMask == 1 {

	} else {
		cp.z = *new(fq).SetUint64(1)
	}

	var err error
	fqX, fqY := new(fq), new(fq)
	_, err = fqX.SetInt(new(big.Int).SetBytes(data[:fqByteLen]))
	if err != nil {
		return err
	}
	cp.x = *fqX
	_, err = fqY.SetInt(new(big.Int).SetBytes(data[fqByteLen:]))
	if err != nil {
		return err
	}
	cp.y = *fqY

	return nil
}

// HashToPoint sets c to the curve point that results from the given slice of bytes
// and returns c. The point is not guaranteed to be in a particular subgroup.
// See https://github.com/Chia-Network/bls-signatures/blob/master/SPEC.md#hashg1
func (c *curvePoint) HashToPoint(buf []byte, ref0 []byte, ref1 []byte) *curvePoint {
	h := blake2b.Sum256(buf)
	sum := blake2b.Sum512(append(h[:], ref0...))
	t0 := new(big.Int)
	g10, _ := new(fq).SetInt(t0.Mod(t0.SetBytes(sum[:]), r))
	sum = blake2b.Sum512(append(h[:], ref1...))
	g11, _ := new(fq).SetInt(t0.Mod(t0.SetBytes(sum[:]), r))

	return c.Add(new(curvePoint).SWEncode(g10), new(curvePoint).SWEncode(g11))
}

// SWUMap maps a value of the finite field to a point in the elliptic curve.
// The point is not guaranteed to be in a particular subgroup.
// SWUMap implements an optimized version of the SWU map.
// See https://eprint.iacr.org/2019/403.pdf - Section 4.4.
func (a *curvePoint) SWUMap(t *fq) *curvePoint {
	n, t0, t1 := new(fq), new(fq), new(fq)
	fqMul(t0, t, t)
	fqMul(t1, t0, t0)
	fqNeg(t0, t0)
	fqAdd(t0, t0, t1)
	fqAdd(n, t0, new(fq).SetUint64(1))
	fqMul(n, n, fqCurveB)

	d := new(fq).Set(t0) // a = -1 for performance
	u := new(fq)
	fqMul(u, n, n)
	fqMul(u, u, n)
	fqMul(t0, d, d)
	fqMul(t1, n, t0)
	fqNeg(t1, t1)
	fqAdd(u, u, t1)
	fqMul(t0, d, t0)
	fqMul(t1, t0, fqCurveB)
	fqAdd(u, u, t1)

	/*
		v := new(fq).Set(t0)

		if t0 == fq{} {
			fqMul(a.x, n, d)
			fqMul(x.y, alpha, v)
			a.z.Set(d)
		} else {
			fqMul(a.x, n, d)
			fqMul(a.y, alpha, v)
			a.z.Set(d)
		}
	*/

	return a
}

// iso11 implements the 11-isogeny from E1´(Fp) to E1(Fp).
// See https://eprint.iacr.org/2019/403.pdf - 4.3 Isogeny maps.
func (a *curvePoint) iso11(b *curvePoint) *curvePoint {
	term := new(fq)
	mul := new(fq)
	var sum [4]fq
	for i, ki := range iso11K {
		sum[i].Set(ki[0])
		mul.SetUint64(1)
		for _, kij := range ki[1:] {
			fqMul(mul, mul, &b.x)
			fqMul(term, kij, mul)
			fqAdd(&sum[i], &sum[i], term)
		}
	}

	fqInv(&sum[1], &sum[1])
	fqMul(&a.x, &sum[0], &sum[1])
	fqInv(&sum[3], &sum[3])
	fqMul(term, &sum[2], &sum[3])
	fqMul(&a.y, &b.y, term)

	return a
}

// SWEncode implements the Shallue and van de Woestijne encoding.
// The point is not guaranteed to be in a particular subgroup.
// See https://www.di.ens.fr/~fouque/pub/latincrypt12.pdf - Algorithm 1.
func (a *curvePoint) SWEncode(b *fq) *curvePoint {
	w, inv := new(fq), new(fq)
	fqMul(w, fqSqrtNegThree, b)
	fqMul(w, w, b)
	fqMul(inv, b, b)
	fqAdd(inv, inv, fqCurveBPlusOne)
	fqInv(inv, inv)
	fqMul(w, w, inv)

	x, y := new(fq), new(fq)
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			fqMul(x, b, w)
			fqSub(x, fqHalfSqrtNegThreeMinusOne, x)
		case 1:
			fqSub(x, &fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}, x)
		case 2:
			fqMul(x, w, w)
			fqInv(x, x)
			fqAdd(x, x, new(fq).SetUint64(1))
		}

		fqMul(y, x, x)
		fqMul(y, y, x)
		fqAdd(y, y, fqCurveB)
		if fqSqrt(y, y) {
			a.x, a.y, a.z = *x, *y, *new(fq).SetUint64(1)
			return a
		}
	}

	return a
}
