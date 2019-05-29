// +build amd64,!generic arm64,!generic

package bls12

//go:noescape
func frAdd(z, x, y *fr)

//go:noescape
func frNeg(z, x *fr)

//go:noescape
func frSub(z, x, y *fr)

//go:noescape
func frMul(z, x, y *fr)
