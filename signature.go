package bls12

import (
	"io"
	"math/big"
)

// PublicKey represents a BLS12-381 public key.
type PublicKey struct {
	*twistPoint
}

// PrivateKey represents a BLS12-381 private key.
type PrivateKey struct {
	PublicKey
	Secret *big.Int
}

// Public returns the public key corresponding to priv.
func (priv *PrivateKey) Public() PublicKey {
	return priv.PublicKey
}

func (priv *PrivateKey) Sign(hash []byte) []byte {
	return Sign(priv, hash)
}

type blsSignature struct {
	*curvePoint
}

func newPrivateKey(index *big.Int) *PrivateKey {
	return &PrivateKey{
		Secret: index,
		PublicKey: PublicKey{
			twistPoint: G2.Element(index),
		},
	}
}

// GenerateKey generates a public and private key pair.
func GenerateKey(reader io.Reader) (*PrivateKey, error) {
	elm, err := randFieldElement(reader)
	if err != nil {
		return nil, err
	}

	return newPrivateKey(elm), nil
}

// Sign signs a hash (which should be the result of hashing a larger message)
// using the private key, priv. If the hash is longer than the bit-length of the
// private key's curve order, the hash will be truncated to that length.
func Sign(priv *PrivateKey, hash []byte) []byte {
	return blsSignature{new(curvePoint).mul(G1.ElementFromHash(hash), priv.Secret)}.marshal()
}

// Verify verifies the signature of hash using the public key(s), pub. Its
// return value records whether the signature is valid.
func Verify(hash []byte, sig []byte, pub ...*PublicKey) bool {
	// TODO
	return false
}

// Aggregate aggregates the signature(s) into a short convincing aggregate signature.
func Aggregate(sig ...[]byte) []byte {
	// TODO
	return []byte{}
}
