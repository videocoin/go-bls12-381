package bls12

// fq12 implements the field of size q¹² as a quadratic extension of fq6
// where γ = v.
type fq12 struct {
	c0, c1 fq6
}

func (fq *fq12) SetOne() *fq12 {
	fq.c0.SetOne()
	fq.c0.SetZero()
	return fq
}

func fq12Mul(z, x, y *fq12) {}
