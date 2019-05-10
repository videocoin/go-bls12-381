package bls12

import (
	"math/big"
)

var (
	// g1X is the x-coordinate of G1's generator in the Montgomery form
	g1X, _ = new(fq).SetString("3685416753713387016781088315183077757961620795782546409894578378688607592378376318836054947676345821548104185464507", Montgomery)

	// g1Y is the y-coordinate of G1's generator in the Montgomery form
	g1Y, _ = new(fq).SetString("1339506544944476473020471379941921221584933875938349620426543736416511423956333506472724655353366534992391756441569", Montgomery)

	g1Gen = &g1Point{curvePoint{*g1X, *g1Y, *fqOne}}

	g10 = []byte("G1_0")

	g11 = []byte("G1_1")

	// g1Cofactor is the cofactor by which to multiply points to map them to G1. (on to the r-torsion). h = (x - 1)2 / 3
	g1Cofactor, _ = bigFromBase10("76329603384216526031706109802092473003")
)

type g1Point struct {
	p curvePoint
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (z *g1Point) ScalarBaseMult(scalar *big.Int) *g1Point {
	return z.ScalarMult(g1Gen, scalar)
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
func (z *g1Point) ScalarMult(x *g1Point, scalar *big.Int) *g1Point {
	z.p.ScalarMult(&x.p, scalar)
	return z
}

// Add returns the sum of (x1,y1) and (x2,y2)
func (z *g1Point) Add(x, y *g1Point) *g1Point {
	z.p.Add(&x.p, &y.p)
	return z
}

// SetBytes uses the Shallue and van de Woestijne encoding.
// The point is guaranteed to be in the subgroup.
func (z *g1Point) SetBytes(buf []byte) *g1Point {
	z.p.SetBytes(buf)
	return z.ScalarMult(z, g1Cofactor)
}

func (z *g1Point) Marshal() []byte {
	return z.p.Marshal()
}

func (z *g1Point) Unmarshal(data []byte) error {
	return z.p.Unmarshal(data)
}
