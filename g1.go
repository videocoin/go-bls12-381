package bls12

import "math/big"

var (
	// G1 is the r-order subgroup of points on the curve
	G1 = new(g1)

	g1Generator = newCurvePoint(g1X, g1Y)
)

type g1 struct{}

func (g1 *g1) Element(index *big.Int) *curvePoint {
	return new(curvePoint).mul(g1Generator, index)
}

func (g1 *g1) ElementFromHash(hash []byte) *curvePoint {
	return &curvePoint{}
}

/*
// hashToCurveSubGroup hashes the msg to a curve sub group point via the given cofactor.
func hashToCurveSubGroup(msg []byte, cofactor *big.Int) *curvePoint {
	point := hashToCurvePoint(msg)
	return point.mul(point, cofactor)
}
*/
