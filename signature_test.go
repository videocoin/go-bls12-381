package bls12

import (
	"crypto/rand"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	// TODO
}

func TestSign(t *testing.T) {
	// TODO
}

func TestVerify(t *testing.T) {
	priv, err := GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pub := priv.Public()
	msg := []byte{7, 8, 9}
	sig := Sign(priv, msg)
	valid, err := Verify(msg, sig, &pub)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Errorf("must be valid")
	}
}
