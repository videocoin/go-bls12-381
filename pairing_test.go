package bls12

import "testing"

func BechmarkPairing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pair(g1Gen, g2Gen)
	}
}
