package bls12

import (
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

func pubKeyFromScalar(scalar *big.Int) *PublicKey {
	return newG2Point().ScalarBaseMult(scalar).ToAffine()
}

func privKeyFromScalar(scalar *big.Int) *PrivateKey {
	return &PrivateKey{
		Secret:    scalar,
		PublicKey: *pubKeyFromScalar(scalar),
	}
}

// GenerateKey generates a public and private key pair.
func GenerateKey(reader io.Reader) (*PrivateKey, error) {
	elem, err := randFieldElement(reader)
	if err != nil {
		return nil, err
	}

	return privKeyFromScalar(elem), nil
}

// Public returns the public key corresponding to priv.
func (priv *PrivateKey) Public() PublicKey {
	return priv.PublicKey
}

// Sign signs a hash (which should be the result of hashing a larger message)
// using the private key, priv.
func Sign(priv *PrivateKey, hash []byte) []byte {
	return newG1Point().SetBytes(hash).ScalarBaseMult(priv.Secret).Marshal()
}

/*
// Verify verifies the signature of hash using the public key, pub. Its
// return value records whether the signature is valid.
// See https://crypto.stanford.edu/~dabo/pubs/papers/BLSmultisig.html -
// BLS signature aggregation - Verify
func Verify(hash []byte, rawSig []byte, pub *PublicKey) (bool, error) {
	sig := newG1Point()
	if err := sig.Unmarshal(rawSig); err != nil {
		return false, err
	}

	return pair(sig, g2Gen).equal(pair(newG1Point().SetBytes(hash), pub)), nil
}
*/

/*
// Aggregate aggregates the signature(s) into a short convincing aggregate signature.
func Aggregate(sig ...[]byte) []byte {
	// TODO
	return []byte{}
}
*/
