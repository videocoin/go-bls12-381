// +build amd64,!generic arm64,!generic

package bls12

//go:noescape
func fqAdd(c, a, b *fq)
