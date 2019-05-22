// +build amd64,!generic arm64,!generic

package bls12

//go:noescape
func fqAdd(z, x, y *fq)

//go:noescape
func fqNeg(z, x *fq)

//go:noescape
func fqSub(z, x, y *fq)

//go:noescape
func fqMul(z, x, y *fq)

//go:noescape
func fqMod(z *fq)
