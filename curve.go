package bls12

import (
	"math/big"

	"golang.org/x/crypto/blake2b"
)

const (
	compressedFormMask  uint8 = 1 << 7
	pointAtInfinityMask uint8 = 1 << 6
)

var curveB, _ = FqFromBase10("4")

// curvePoint is an elliptic curve point in projective coordinates.
// The elliptic curve is defined by the following equation y²=x³+3.
type curvePoint struct {
	x, y, z Fq
}

func newCurvePoint(x, y Fq) *curvePoint {
	return &curvePoint{
		x: x,
		y: y,
		z: Fq0,
	}
}

// curvePointFromHash converts the hash to a curve point.
// The point is not guaranteed to be in a particular subgroup.
func curvePointFromHash(hash []byte) *curvePoint {
	// See https://github.com/Chia-Network/bls-signatures/blob/master/SPEC.md#hashg1
	h256, _ := blake2b.New256(nil)
	h256.Write(hash)
	h := h256.Sum(nil)

	h512, _ := blake2b.New512(nil)
	h512.Write(h)
	h512.Write([]byte("G1_0"))
	t0 := curvePointFromFq(fqFromHash(h512.Sum(nil)))

	h512.Reset()
	h512.Write(h)
	h512.Write([]byte("G1_1"))
	t1 := curvePointFromFq(fqFromHash(h512.Sum(nil)))

	return new(curvePoint).add(t0, t1)
}

func (cp *curvePoint) add(a, b *curvePoint) *curvePoint {
	// See https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	z1z1, z2z2 := new(Fq), new(Fq)
	fqSqr(z1z1, &a.z)
	fqSqr(z2z2, &b.z)

	u1, u2 := new(Fq), new(Fq)
	fqMul(u1, &a.x, z2z2)
	fqMul(u2, &b.x, z1z1)

	s1, s2 := new(Fq), new(Fq)
	fqMul(s1, &a.y, &b.z)
	fqMul(s1, s1, z2z2)
	fqMul(s2, &b.y, &a.z)
	fqMul(s2, s2, z1z1)

	h, i, j, r, v := new(Fq), new(Fq), new(Fq), new(Fq), new(Fq)
	fqSub(h, u2, u1)
	fqDbl(i, h)
	fqSqr(i, i)
	fqMul(j, h, i)
	fqSub(r, s2, s1)
	fqDbl(r, r)
	fqMul(v, u1, i)

	t0, t1 := new(Fq), new(Fq)
	fqDbl(t0, v)
	fqSqr(&cp.x, r)
	fqSub(&cp.x, &cp.x, j)
	fqSub(&cp.x, &cp.x, t0)
	fqDbl(t0, s1)
	fqMul(t0, t0, j)
	fqMul(t1, r, &cp.x)
	fqMul(&cp.y, r, v)
	fqSub(&cp.y, &cp.y, t1)
	fqSub(&cp.y, &cp.y, t0)
	fqMul(&cp.z, &a.z, &b.z)
	fqSqr(&cp.z, &cp.z)
	fqSub(&cp.z, &cp.z, z1z1)
	fqSub(&cp.z, &cp.z, z2z2)
	fqMul(&cp.z, &cp.z, h)

	return cp
}

func (cp *curvePoint) double(p *curvePoint) *curvePoint {
	// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
	a, b, c, d, e, f, t0 := new(Fq), new(Fq), new(Fq), new(Fq), new(Fq), new(Fq), new(Fq)

	fqSqr(a, &p.x)
	fqSqr(b, &p.y)
	fqSqr(c, b)
	fqAdd(d, &p.x, b)
	fqSqr(d, d)
	fqSub(d, d, a)
	fqSub(d, d, c)
	fqDbl(d, d)
	fqDbl(e, a)
	fqAdd(e, e, a)
	fqSqr(f, e)

	fqDbl(&cp.x, d)
	fqSub(&cp.x, f, &cp.x)
	fqAdd(t0, c, c)
	fqAdd(t0, t0, t0)
	fqDbl(t0, t0)
	fqSub(&cp.y, d, &cp.x)
	fqMul(&cp.y, e, &cp.y)
	fqSub(&cp.y, &cp.y, t0)
	fqMul(&cp.z, &cp.y, &cp.z)
	fqAdd(&cp.z, &cp.z, &cp.z)

	return cp
}

func (cp *curvePoint) mul(p *curvePoint, scalar *big.Int) *curvePoint {
	// See https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add
	q := new(curvePoint)
	for i := scalar.BitLen(); i > 0; i-- {
		q.double(q)
		if scalar.Bit(i) == 1 {
			q.add(q, p)
		}
	}

	return q
}

// isInfinty check if the point is a point at "infinity"
func (cp *curvePoint) isInfinity() bool {
	return cp.z == Fq0
}

func (cp *curvePoint) makeAffine() {
	// TODO create new curve point
	if cp.isInfinity() {
		//  If this bit is set, the remaining bits of the group element's encoding should be set to zero.
		//pointAtInfinityMask
	}
	zInv := new(Fq)
	fqInv(zInv, &cp.z)
	fqMul(&cp.x, &cp.x, zInv)
	fqMul(&cp.y, &cp.y, zInv)
}

// Marshal converts a point into the compressed form specified in
// https://github.com/zkcrypto/pairing/tree/master/src/bls12_381#serialization.
func (cp *curvePoint) marshal() []byte {
	cp.makeAffine()
	// See https://github.com/zkcrypto/pairing/tree/master/src/bls12_381#serialization
	// TODO: https://golang.org/src/crypto/elliptic/elliptic.go?s=8258:8305#L296
	//cp.MakeAffine()

	var x Fq
	montgomeryDecode(&x, &cp.x)

	xBytes := x.Bytes()
	xBytes[0] &= compressedFormMask

	return xBytes

}

// unmarshalCurvePoint converts a point, serialized by Marshal, into an x, y pair.
// It is an error if the point is not in compressed form or is not on the curve.
// On error, x = nil.
func unmarshalCurvePoint(data []byte) (*curvePoint, error) {
	if len(data) != fqCompressedLen {
		// TODO (error)
		return nil, nil
	}
	if data[0]&compressedFormMask == 0 {
		// TODO (error)
		return nil, nil
	}

	return &curvePoint{}, nil
}

func curvePointFromFq(elm Fq) *curvePoint {
	return newCurvePoint(coordinatesFromFq(elm))
}

// hashToCurveSubGroup hashes the msg to a specific curve subgroup.
// cofactor https://crypto.stackexchange.com/questions/33028/order-and-cofactor-of-the-base-point
func hashToCurveSubGroup(msg []byte, cofactor *big.Int) *curvePoint {
	return new(curvePoint).mul(curvePointFromHash(msg), cofactor)
}
