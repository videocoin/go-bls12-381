// +build amd64,!generic arm64,!generic

package bls12

//go:noescape
func FqAdd(z, x, y *Fq)

func FqDbl(z, x *Fq) { FqAdd(z, x, x) }

//go:noescape
func FqNeg(z, x *Fq)

//go:noescape
func FqSub(z, x, y *Fq)

//go:noescape
func FqBasicMul(z *FqLarge, x, y *Fq)

//go:noescape
func FqREDC(z *Fq, x *FqLarge)

//go:noescape
func FqMul(z, x, y *Fq)

func FqSqr(z, x *Fq) { FqMul(z, x, x) }

//go:noescape
func FqSqrt(z, x *Fq) bool

//go:noescape
func FqCube(z, x *Fq)

//go:noescape
func FqExp(z, x *Fq, y []uint64)

func FqInv(c, x *Fq) { FqExp(c, x, qMinus2[:]) }
