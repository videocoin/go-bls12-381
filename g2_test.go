package bls12

import (
	"crypto/rand"
	"testing"
)

func BenchmarkG2(b *testing.B) {
	x, _ := randFieldElement(rand.Reader)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		new(g2Point).ScalarBaseMult(x)
	}
}
