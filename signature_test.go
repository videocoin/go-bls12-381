package bls12

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	// TODO
}

func TestSignAndVerify(t *testing.T) {
	// TODO
	/*
		priv, _ := GenerateKey(rand.Reader)

		hashed := []byte("testing")
		sig := Sign(priv, hashed)

		if !Verify(&priv.PublicKey, hashed, sig) {
			t.Errorf("%s: Verify failed", tag)
		}

		hashed[0] ^= 0xff
		if Verify(&priv.PublicKey, hashed, r, s) {
			t.Errorf("%s: Verify always works!", tag)
		}
	*/
}

func TestZeroHashSignature(t *testing.T) {
	// TODO
	/*
		zeroHash := make([]byte, 64)
		privKey, err := GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}

		// Sign a hash consisting of all zeros.
		r, s, err := Sign(rand.Reader, privKey, zeroHash)
		if err != nil {
			panic(err)
		}

		// Confirm that it can be verified.
		if !Verify(&privKey.PublicKey, zeroHash, r, s) {
			t.Errorf("zero hash signature verify failed for %T", curve)
		}
	*/
}
