// +build amd64,!generic arm64,!generic

package bls12

//go:noescape
func fqAdd(z, x, y *Fq)

func fqDbl(z, x *Fq) { fqAdd(z, x, x) }

//go:noescape
func fqNeg(z, x *Fq)

//go:noescape
func fqSub(z, x, y *Fq)

//go:noescape
func fqBasicMul(z *FqLarge, x, y *Fq)

//go:noescape
func fqREDC(z *Fq, x *FqLarge)

//go:noescape
func fqMul(z, x, y *Fq)

func fqSqr(z, x *Fq) { fqMul(z, x, x) }

//go:noescape
func fqSqrt(z, x *Fq) bool

//go:noescape
func fqCube(z, x *Fq)

//go:noescape
func fqExp(z, x *Fq, y []uint64)

func fqInv(c, x *Fq) { fqExp(c, x, qm2[:]) }
