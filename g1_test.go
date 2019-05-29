package bls12

import (
	"crypto/rand"
	"testing"
)

func TestG1PointScalarBaseMult(t *testing.T) {
	// TODO
}

func TestG1PointScalarMult(t *testing.T) {
	// TODO
}

func TestG1PointAdd(t *testing.T) {
	// TODO
}

func TestG1PointSetBytes(t *testing.T) {
	// TODO
}

func TestG1PointMarshal(t *testing.T) {
	// TODO
}

func TestG1PointUnmarshal(t *testing.T) {
	// TODO
}

func BenchmarkG1(b *testing.B) {
	x, _ := randFieldElement(rand.Reader)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		new(G1Point).ScalarBaseMult(x)
	}
}
