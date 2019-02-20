// +build amd64,!generic arm64,!generic

package bls12

//go:noescape
func fqAdd(c, a, b *fq)

//go:noescape
func fqMod(a *fq, head uint64)
