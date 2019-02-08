package bls12

var (
	curveB = newFq("4")

	g1Generator = &curvePoint{
		x:          *newFq(g1X),
		y:          *newFq(g1Y),
		z:          *newFq("1"),
		isInfinity: false, // TODO
	}
)

// curvePoint implements the elliptic curve y²=x³+4.
type curvePoint struct {
	x, y, z    fq
	isInfinity bool
}

func (c *curvePoint) IsInfinity() bool {
	return c.isInfinity
}

// Add sets c to the sum a+b and returns c.
func (c *curvePoint) Add(a, b *curvePoint) *curvePoint {
	/*
		if a == nil {
			*a = *c
		}

		if a.IsInfinity() {
			c.Set(b)
			return c
		}
		if b.IsInfinity() {
			c.Set(a)
			return c
		}

		// See https://hyperelliptic.org/EFD/g1p/auto-code/shortw/jacobian-3/addition/mmadd-2007-bl.op3

	*/
	return c
}

func (c *curvePoint) Double(a *curvePoint) *curvePoint {
	if a == nil {
		*a = *c
	}
	// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l

	return c
}
