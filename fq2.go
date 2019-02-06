package bls

// Fq2 is an element of Fq2, represented by c0 + c1 * u.
type Fq2 struct {
	C0 *Fq
	C1 *Fq
}
