// +build amd64,!generic arm64,!generic

package bls12

//go:noescape
func fqAdd(z, x, y *fq)

//go:noescape
func fqNeg(z, x *fq)

//go:noescape
func fqSub(z, x, y *fq)

//go:noescape
func fqBasicMul(z *fqLarge, x, y *fq)

//go:noescape
func fqREDC(z *fq, x *fqLarge)

//go:noescape
func fqMul(z, x, y *fq)

//go:noescape
func fqSqrt(z, x *fq) bool

//go:noescape
func fqExp(z, x *fq, y []uint64)

func fqInv(c, x *fq) { fqExp(c, x, qMinus2[:]) }

func fqLargeSub(c, a, b *fqLarge)

func fqLargeAdd(c, a, b *fqLarge)
