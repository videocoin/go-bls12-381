package bls12

import (
	"crypto/rand"
	"io"
	"math/big"
)

// PublicKey represents a BLS public key.
type PublicKey = g2Point

// PrivateKey represents a BLS private key.
type PrivateKey struct {
	PublicKey
	Secret *big.Int
}

func privKeyFromScalar(k *big.Int) *PrivateKey {
	priv := new(PrivateKey)
	priv.Secret = k
	priv.PublicKey = *new(g2Point).ScalarBaseMult(k).ToAffine()
	return priv
}

// randInt returns a random scalar between 0 and max.
func randInt(reader io.Reader, max *big.Int) (n *big.Int, err error) {
	for {
		n, err = rand.Int(reader, max)
		if n.Sign() > 0 || err != nil {
			return
		}
	}
}

// randFieldElement returns a random scalar between 0 and r.
func randFieldElement(reader io.Reader) (*big.Int, error) {
	return randInt(reader, r)
}

// GenerateKey generates a public and private key pair.
func GenerateKey(reader io.Reader) (*PrivateKey, error) {
	k, err := randFieldElement(reader)
	if err != nil {
		return nil, err
	}

	return privKeyFromScalar(k), nil
}

// Public returns the public key corresponding to priv.
func (priv *PrivateKey) Public() PublicKey {
	return priv.PublicKey
}

// Sign signs a hash (which should be the result of hashing a larger message)
// using the private key, priv.
func Sign(priv *PrivateKey, hash []byte) []byte {
	return new(g1Point).ScalarMult(new(g1Point).SetBytes(hash), priv.Secret).Marshal()
}

// Verify verifies the signature of hash using the public key, pub. Its
// return value records whether the signature is valid.
func Verify(hash []byte, rawSig []byte, pub *PublicKey) (bool, error) {
	var sig g1Point
	if err := sig.Unmarshal(rawSig); err != nil {
		return false, err
	}
	return Pair(&sig, g2Gen).Equal(Pair(new(g1Point).SetBytes(hash), pub)), nil
}
