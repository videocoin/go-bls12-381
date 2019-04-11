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
	// URGENT fix me support negative numbers
	fqMontNeg1, _ = fqMontgomeryFromBase10("1")
)

// curvePoint is an elliptic curve point in projective coordinates.
// The elliptic curve is defined by the following equation y²=x³+3.
type curvePoint struct {
	x, y, z fq
}

func newCurvePoint(x, y fq) *curvePoint {
	return &curvePoint{
		x: x,
		y: y,
		z: fq0,
	}
}

func (cp *curvePoint) Set(p *curvePoint) *curvePoint {
	cp.x, cp.y, cp.z = p.x, p.y, p.z

	return cp
}

func (cp *curvePoint) IsInfinity() bool {
	return cp.z == fq0
}

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
	fqSqr(z1z1, &a.z)
	fqSqr(z2z2, &b.z)

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
	fqDbl(i, h)
	fqSqr(i, i)
	fqMul(j, h, i)
	fqSub(r, s2, s1)
	fqDbl(r, r)
	fqMul(v, u1, i)

	x, y, z, t0, t1 := new(fq), new(fq), new(fq), new(fq), new(fq)
	fqDbl(t0, v)
	fqAdd(t0, t0, j)
	fqSqr(x, r)
	fqSub(x, x, t0)

	fqDbl(t0, s1)
	fqMul(t0, t0, j)
	fqSub(t1, v, x)
	fqMul(t1, t1, r)
	fqSub(y, t1, t0)

	fqAdd(z, &a.z, &b.z)
	fqSqr(z, z)
	fqAdd(t0, z1z1, z2z2)
	fqSub(z, z, t0)
	fqMul(z, z, h)

	cp.x, cp.y, cp.z = *x, *y, *z

	return cp
}

func (cp *curvePoint) Double(p *curvePoint) *curvePoint {
	// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
	a, b, c, d, e, f := new(fq), new(fq), new(fq), new(fq), new(fq), new(fq)
	fqSqr(a, &p.x)
	fqSqr(b, &p.y)
	fqSqr(c, b)
	fqMul(d, &p.x, b)
	fqDbl(d, d)
	fqDbl(d, d)
	fqDbl(e, a)
	fqAdd(e, e, a)
	fqSqr(f, e)

	x, y, z, t0 := new(fq), new(fq), new(fq), new(fq)
	fqDbl(x, d)
	fqSub(x, f, x)
	fqDbl(t0, c)
	fqDbl(t0, t0)
	fqDbl(t0, t0)
	fqSub(y, d, x)
	fqMul(y, y, e)
	fqSub(y, y, t0)
	fqMul(z, &p.y, &p.z)
	fqDbl(z, z)

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
	return cp.z == fq0
}

func (cp *curvePoint) ToAffine() {
	// TODO create new curve point
	if cp.isInfinity() {
		//  If this bit is set, the remaining bits of the group element's encoding should be set to zero.
		//pointAtInfinityMask
	}
	zInv := new(fq)
	fqInv(zInv, &cp.z)
	fqMul(&cp.x, &cp.x, zInv)
	fqMul(&cp.y, &cp.y, zInv)
}

// CompressedEncode converts a curve point into the uncompressed form specified in
// See https://github.com/zkcrypto/pairing/tree/master/src/bls12_381#serialization.
func (cp *curvePoint) Marshal() []byte {
	cp.ToAffine()
	x, y := new(fq), new(fq)
	montgomeryDecode(x, &cp.x)
	montgomeryDecode(y, &cp.y)

	ret := make([]byte, fqByteLen*2)

	xBytes := x.Bytes()
	copy(ret, xBytes)
	yBytes := y.Bytes()
	copy(ret[fqByteLen:], yBytes)

	// TODO review
	if cp.IsInfinity() {
		ret[0] |= pointAtInfinityMask
	}

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
	fqSqr(inv, t)
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
			fqSqr(x, w)
			fqInv(x, x)
			fqAdd(x, x, &fqMont1)
		}

		fqCube(y, x)
		fqAdd(y, y, &fqMontCurveB)
		if fqSqrt(y, y) {
			cp.x = *x
			cp.x = *y
			return cp
		}
	}

	return cp
}
