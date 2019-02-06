package bls

// Fq12 is an element of Fq12, represented by c0 + c1 * w.
type Fq12 struct {
	C0 *Fq6
	C1 *Fq6
}
