package bls12

import "math/big"

type g1Point = curvePoint

var (
	// G1 is the r-order subgroup of points on the curve
	G1 = new(g1)

	g1Generator = newCurvePoint(g1X, g1Y)

	unmarshalG1Point = unmarshalCurvePoint
)

type g1 struct{}

func (g1 *g1) Element(index *big.Int) *g1Point {
	return new(g1Point).mul(g1Generator, index)
}

func (g1 *g1) ElementFromHash(hash []byte) *g1Point {
	return hashToCurveSubGroup(hash, g1Cofactor)
}
