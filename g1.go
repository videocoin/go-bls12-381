package bls12

var (
	g1Gen = &curvePoint{
		x: g1X,
		y: g1Y,
	}
)

// G1 is an abstract cyclic group. The zero value is suitable for use as the
// output of an operation, but cannot be used as an input.
type G1 struct{}
