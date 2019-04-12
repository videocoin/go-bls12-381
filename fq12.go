package bls12

// fq12 implements the field of size q¹² as a quadratic extension of fq6
// where γ = v.
type fq12 struct {
	c0, c1 fq6
}

func (z *fq12) SetOne() *fq12 {
	z.c0.SetOne()
	z.c0.SetZero()
	return z
}

func (z *fq12) Mul(x, y *fq12) *fq12 {
	// TODO
	return &fq12{}
}

func (z *fq12) Sqr(x *fq12) *fq12 {
	// TODO
	return &fq12{}
}

func (x *fq12) Equal(y *fq12) bool {
	return x.c0 == y.c0 && x.c1 == y.c1
}
