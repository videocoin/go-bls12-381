package bls12

import (
	"crypto/rand"
	"io"
)

// PrivateKey represents a BLS12-381 private key.
type PrivateKey struct {
	PublicKey
}

// PublicKey represents a BLS12-381 public key.
type PublicKey struct{}

// Public returns the public key corresponding to priv.
func (priv *PrivateKey) Public() PublicKey {
	return PublicKey{}
}

// Sign signs a hash with priv.
func (priv *PrivateKey) Sign(rand io.Reader, hash []byte) []byte {
	return []byte{}
}

// Sign signs a hash (which should be the result of hashing a larger message)
// using the private key, priv. If the hash is longer than the bit-length of the
// private key's curve order, the hash will be truncated to that length. The security
// of the private key depends on the entropy of rand.
func Sign(rand io.Reader, priv *PrivateKey, hash []byte) []byte {
	return []byte{}
}

// GenerateKey generates a public and private key pair.
func GenerateKey(seed []byte) (*PrivateKey, error) {
	if seed == nil {
		_, err := rand.Read(seed)
		if err != nil {
			return nil, err
		}
	}

	return &PrivateKey{}, nil
}

// Verify verifies the signature of hash using the public key(s), pub. Its
// return value records whether the signature is valid.
func Verify(hash []byte, sig []byte, pub ...*PublicKey) bool {
	return false
}

// Aggregate aggregates the signature(s) into a short convincing aggregate signature.
func Aggregate(sig ...[]byte) []byte {
	return []byte{}
}
