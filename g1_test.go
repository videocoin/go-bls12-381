package bls12

import (
	"crypto/rand"
	"testing"
)

func TestG1PointUnmarshal(t *testing.T) {
	// TODO
}

func BenchmarkG1(b *testing.B) {
	x, _ := randFieldElement(rand.Reader)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		new(g1Point).ScalarBaseMult(x)
	}
}
