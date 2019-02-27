package bls12

// fq6 is an element of Fq6 = Fq²[Y]/(Y³ − γ), where γ
// is a quadratic non-residue in Fq and γ = √β is a
// cubic non-residue in Fq² with a value of X + 1.
// See https://eprint.iacr.org/2006/471.pdf - "6.2 Cubic over quadratic"
type fq6 struct {
	c0, c1, c2 fq2
}
