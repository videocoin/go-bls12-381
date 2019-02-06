package bls

// Fq6 is an element of Fq6, represented by c0 + c1 * v + c2 * v^(2).
type Fq6 struct {
	C0 *Fq2
	C1 *Fq2
	C2 *Fq2
}
