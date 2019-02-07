package bls12

var (
	curveB = newFq("4")

	g1Generator = &curvePoint{
		x: *newFq(g1X),
		y: *newFq(g1Y),
		z: *newFq("1"),
	}
)

// curvePoint implements the elliptic curve y²=x³+4.
type curvePoint struct {
	x, y, z fq
}
