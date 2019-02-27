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

func fqSqr(z, x *fq) {
	fqMul(z, x, x)
}

func fqDbl(z, x *fq) {
	fqAdd(z, x, x)
}
