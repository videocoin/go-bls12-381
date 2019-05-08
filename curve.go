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
	fqMontCurveB, _             = fqMontgomeryFromBase10("4")
	fqMontCurveBPlus1, _        = fqMontgomeryFromBase10("5")
	fqMontSqrtNeg3, _           = fqMontgomeryFromBase10("1586958781458431025242759403266842894121773480562120986020912974854563298150952611241517463240701")
	fqMontHalfSqrtNeg3Minus1, _ = fqMontgomeryFromBase10("793479390729215512621379701633421447060886740281060493010456487427281649075476305620758731620350")
	// fixme:  add support negative numbers
	fqMontNeg1, _ = fqMontgomeryFromBase10("1")
)

// curvePoint is an elliptic curve point in projective coordinates.
// The elliptic curve is defined by the following equation y²=x³+3.
type curvePoint struct {
	x, y, z fq
}

// Set sets cp to the value of p and returns cp.
func (cp *curvePoint) Set(p *curvePoint) *curvePoint {
	cp.x, cp.y, cp.z = p.x, p.y, p.z
	return cp
}

func (cp *curvePoint) IsInfinity() bool {
	return cp.z == fqZero
}

// Add sets cp to the sum a+b and returns cp.
func (cp *curvePoint) Add(a, b *curvePoint) *curvePoint {
	// TODO is infinity confirm
	if a.IsInfinity() {
		return cp.Set(b)
	}

	// TODO is infinity confirm
	if b.IsInfinity() {
		return cp.Set(a)
	}

	if a.Equal(b) {
		return cp.Double(a) // faster than Add
	}

	// See https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	z1z1, z2z2 := new(fq), new(fq)
	fqMul(z1z1, &a.z, &a.z)
	fqMul(z2z2, &a.z, &a.z)

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

	x, y, z, t0, t1 := new(fq), new(fq), new(fq), new(fq), new(fq)
	fqAdd(t0, v, v)
	fqAdd(t0, t0, j)
	fqMul(x, r, r)
	fqSub(x, x, t0)

	fqAdd(t0, s1, s1)
	fqMul(t0, t0, j)
	fqSub(t1, v, x)
	fqMul(t1, t1, r)
	fqSub(y, t1, t0)

	fqAdd(z, &a.z, &b.z)
	fqMul(z, z, z)
	fqAdd(t0, z1z1, z2z2)
	fqSub(z, z, t0)
	fqMul(z, z, h)

	cp.x, cp.y, cp.z = *x, *y, *z

	return cp
}

// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
func (cp *curvePoint) Double(p *curvePoint) *curvePoint {
	a, b, c, d, e, f := new(fq), new(fq), new(fq), new(fq), new(fq), new(fq)
	fqMul(a, &p.x, &p.x)
	fqMul(b, &p.y, &p.y)
	fqMul(c, b, b)
	fqMul(d, &p.x, b)
	fqAdd(d, d, d)
	fqAdd(d, d, d)
	fqAdd(e, a, a)
	fqAdd(e, e, a)
	fqMul(f, e, e)

	x, y, z, t0 := new(fq), new(fq), new(fq), new(fq)
	fqAdd(x, d, d)
	fqSub(x, f, x)
	fqAdd(t0, c, c)
	fqAdd(t0, t0, t0)
	fqAdd(t0, t0, t0)
	fqSub(y, d, x)
	fqMul(y, y, e)
	fqSub(y, y, t0)
	fqMul(z, &p.y, &p.z)
	fqAdd(z, z, z)

	cp.x, cp.y, cp.z = *x, *y, *z

	return cp
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
// See https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add
func (cp *curvePoint) ScalarMult(p *curvePoint, scalar *big.Int) *curvePoint {
	q := new(curvePoint)
	for i := scalar.BitLen() - 1; i >= 0; i-- {
		q.Double(q)
		if scalar.Bit(i) == 1 {
			q.Add(q, p)
		}
	}

	return cp.Set(q)
}

func (cp *curvePoint) Equal(p *curvePoint) bool {
	return cp.x == p.x && cp.y == p.y && cp.z == p.z
}

// isInfinty check if the point is a point at "infinity"
func (cp *curvePoint) isInfinity() bool {
	return cp.z == fqZero
}

func (cp *curvePoint) ToAffine() *curvePoint {
	if cp.z.IsOne() {
		return cp
	}

	// TODO create new curve point
	if cp.isInfinity() {
		//  If this bit is set, the remaining bits of the group element's encoding should be set to zero.
		//pointAtInfinityMask
	}

	zInv := new(fq)
	fqInv(zInv, &cp.z)
	fqMul(&cp.x, &cp.x, zInv)
	fqMul(&cp.y, &cp.y, zInv)
	cp.z = fqMontOne

	return cp
}

// CompressedEncode converts a curve point into the uncompressed form specified in
// See https://github.com/zkcrypto/pairing/tree/master/src/bls12_381#serialization.
func (cp *curvePoint) Marshal() []byte {
	cp.ToAffine()

	x := new(fq).MontgomeryDecode(&cp.x)
	y := new(fq).MontgomeryDecode(&cp.y)

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
		cp.z = fqMontOne
	}

	var err error
	cp.x, err = fqMontgomeryFromBig(new(big.Int).SetBytes(data[:fqByteLen]))
	if err != nil {
		return err
	}
	cp.y, err = fqMontgomeryFromBig(new(big.Int).SetBytes(data[fqByteLen:]))
	if err != nil {
		return err
	}

	return nil
}

/*
// unmarshalCurvePoint converts a point, serialized by Marshal, into an x, y pair.
// It is an error if the point is not in compressed form or is not on the curve.
// On error, x = nil.
func unmarshalCurvePoint(data []byte) (*curvePoint, error) {


	return &curvePoint{}, nil
}
*/

// TODO desc
// The point is not guaranteed to be in a particular subgroup.
// See https://github.com/Chia-Network/bls-signatures/blob/master/SPEC.md#hashg1
func (cp *curvePoint) SetBytes(buf []byte) *curvePoint {
	h := blake2b.Sum256(buf)
	sum := blake2b.Sum512(append(h[:], g10...))
	g10, _ := fqMontgomeryFromBig(new(big.Int).Mod(new(big.Int).SetBytes(sum[:]), q))
	sum = blake2b.Sum512(append(h[:], g11...))
	g11, _ := fqMontgomeryFromBig(new(big.Int).Mod(new(big.Int).SetBytes(sum[:]), q))

	return cp.Add(new(curvePoint).SWEncode(&g10), new(curvePoint).SWEncode(&g11))
}

// SWEncode implements the Shallue and van de Woestijne encoding.
// The point is not guaranteed to be in a particular subgroup.
// See https://www.di.ens.fr/~fouque/pub/latincrypt12.pdf - Algorithm 1
func (cp *curvePoint) SWEncode(t *fq) *curvePoint {
	w, inv := new(fq), new(fq)
	fqMul(w, &fqMontSqrtNeg3, t)
	fqMul(w, w, t)
	fqMul(inv, t, t)
	fqAdd(inv, inv, &fqMontCurveBPlus1)
	fqInv(inv, inv)
	fqMul(w, w, inv)

	x, y := new(fq), new(fq)
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			fqMul(x, t, w)
			fqSub(x, &fqMontHalfSqrtNeg3Minus1, x)
		case 1:
			fqSub(x, &fqMontNeg1, x)
		case 2:
			fqMul(x, w, w)
			fqInv(x, x)
			fqAdd(x, x, &fqMontOne)
		}

		fqMul(y, x, x)
		fqMul(y, y, x)
		fqAdd(y, y, &fqMontCurveB)
		if fqSqrt(y, y) {
			cp.x, cp.y, cp.z = *x, *y, fqMontOne
			return cp
		}
	}

	return cp
}
