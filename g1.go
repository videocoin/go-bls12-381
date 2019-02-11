package bls12

var (
	g1Generator = &curvePoint{
		x: *newFq(g1X),
		y: *newFq(g1Y),
		z: *newFq("1"),
	}
)

// G1 is an abstract cyclic group. The zero value is suitable for use as the
// output of an operation, but cannot be used as an input.
type G1 struct{}
