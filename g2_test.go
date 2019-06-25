package bls12

import (
	"crypto/rand"
	"testing"
)

func TestG2PointSet(t *testing.T) {
	// TODO
}

func TestG2PointEqual(t *testing.T) {
	// TODO
}

func TestG2PointAdd(t *testing.T) {
	// TODO
}

func TestG2PointBaseScalarMult(t *testing.T) {
	// TODO
}

func TestG2PointScalarMult(t *testing.T) {
	// TODO
}

func TestG2PointToAffine(t *testing.T) {
	// TODO
}

func BenchmarkG2(b *testing.B) {
	x, _ := RandFieldElement(rand.Reader)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		new(G2Point).ScalarBaseMult(x)
	}
}
