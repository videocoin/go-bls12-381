package bls12

import "math/big"

var (
	g1Generator = newCurvePoint(g1X, g1Y)

	G1 = new(g1)
)

// g1 is the r-order subgroup of points on the curve
type g1 struct{}

func (g1 *g1) element(index *big.Int) *twistPoint {
	return new(twistPoint).mul(g2Generator, index)
}

func curvePointFromFq(elm fq) *curvePoint {
	return newCurvePoint(coordinatesFromFq(elm))
}

/*
// hashToCurveSubGroup hashes the msg to a curve sub group point via the given cofactor.
func hashToCurveSubGroup(msg []byte, cofactor *big.Int) *curvePoint {
	point := hashToCurvePoint(msg)
	return point.mul(point, cofactor)
}
*/
