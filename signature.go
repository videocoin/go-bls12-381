package bls12

import (
	"crypto"
	"io"
	"math/big"
)

// PublicKey represents a BLS12-381 public key.
type PublicKey = g2Point

// PrivateKey represents a BLS12-381 private key.
type PrivateKey struct {
	PublicKey
	Secret *big.Int
}

type blsSignature = g1Point

// PrivateKeyFromScalar returns a new private key instance.
func PrivateKeyFromScalar(scalar *big.Int) *PrivateKey {
	return &PrivateKey{
		Secret:    scalar,
		PublicKey: *G2.ScalarBaseMult(scalar),
	}
}

// GenerateKey generates a public and private key pair.
func GenerateKey(reader io.Reader) (*PrivateKey, error) {
	elem, err := randFieldElement(reader)
	if err != nil {
		return nil, err
	}

	return PrivateKeyFromScalar(elem), nil
}

// Public returns the public key corresponding to priv.
func (priv *PrivateKey) Public() crypto.PublicKey {
	return &priv.PublicKey
}

/*
// Sign signs a hash (which should be the result of hashing a larger message)
// using the private key, priv. If the hash is longer than the bit-length of the
// private key's curve order, the hash will be truncated to that length.

func Sign(priv *PrivateKey, hash []byte) []byte {
	return blsSignature{new(curvePoint).mul(G1.ElementFromHash(hash), priv)}.marshal()
}


// Verify verifies the signature of hash using the public key(s), pub. Its
// return value records whether the signature is valid.
func Verify(hash []byte, sig []byte, pubKey *PublicKey) (bool, error) {
	sigPoint, err := unmarshalSignature(sig)
	if err != nil {
		return false, err
	}

	return pair(sigPoint, g2Generator).equal(pair(G1.ElementFromHash(hash), pubKey)), nil
}


// Aggregate aggregates the signature(s) into a short convincing aggregate signature.
func Aggregate(sig ...[]byte) []byte {
	// TODO
	return []byte{}
}
*/
