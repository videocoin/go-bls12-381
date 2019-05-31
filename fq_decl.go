// +build amd64,!generic arm64,!generic

package bls12

// fqAdd sets z to the sum x+y.
func fqAdd(z, x, y *fq)

// fqNeg sets z to -x.
func fqNeg(z, x *fq)

// fqSub sets z to the difference x-y.
func fqSub(z, x, y *fq)

// fqMul sets z to the product x*y.
func fqMul(z, x, y *fq)
