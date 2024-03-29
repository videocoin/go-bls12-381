package sig1

import (
	"crypto/rand"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	priv, err := GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	hashed := []byte("testing")

	sig := Sign(priv, hashed)

	if valid := Verify(hashed, sig, &priv.PublicKey); !valid {
		t.Errorf("Verify failed")
	}

	hashed[0] ^= 0xff
	if valid := Verify(hashed, sig, &priv.PublicKey); valid {
		t.Errorf("Verify always works!")
	}
}

func TestZeroHashSignature(t *testing.T) {
	priv, err := GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	zeroHash := make([]byte, 64)
	sig := Sign(priv, zeroHash)

	// Confirm that it can be verified.
	if valid := Verify(zeroHash, sig, &priv.PublicKey); !valid {
		t.Errorf("zero hash signature verify failed")
	}
}
