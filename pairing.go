package bls12

/*
// gt is an element of the order q multiplicative subgroup of Fq12.
type gt struct {
	elem fq12
}

func (z *gt) Pair(x *g1Point, y *g2Point) {

}

func (x *gt) Equal(y *gtElem) bool {
	return x.elem == y.elem
}

// computes (p^12 − 1)/r
func finalExp()

//  maps two points into a element of a finite field
func miller(a *curvePoint, b *twistPoint) *fq12 {}



func miller(a *curvePoint, b *twistPoint) *fq12 {
	f := fq12{1}
}
*/

func miller(a *curvePoint, b *twistPoint) *fq12 {
	f := new(fq12).SetOne()

	return f
}
