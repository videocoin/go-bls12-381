package bls12

type gtPoint fq12

func (a *gtPoint) equal(b *gtPoint) bool {
	return true
}

func pair(p1 *g1Point, p2 *g2Point) *gtPoint {
	return &gtPoint{}
}
