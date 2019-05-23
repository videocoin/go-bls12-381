package bls12

import (
	"math/big"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	// TODO
}

func TestSignAndVerify(t *testing.T) {
	//priv, _ := GenerateKey(rand.Reader)
	priv := privKeyFromScalar(new(big.Int).SetUint64(15))
	hashed := []byte("testing")
	sig := Sign(priv, hashed)

	if valid, _ := Verify(hashed, sig, &priv.PublicKey); !valid {
		t.Errorf("Verify failed")
		return
	}

	hashed[0] ^= 0xff
	if valid, _ := Verify(hashed, sig, &priv.PublicKey); valid {
		t.Errorf("Verify always works!")
		return
	}
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
