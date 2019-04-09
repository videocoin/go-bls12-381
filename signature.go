package bls12

import (
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

// privateKeyFromScalar returns a new private key instance.
func privKeyFromScalar(scalar *big.Int) *PrivateKey {
	return &PrivateKey{
		Secret:    scalar,
		PublicKey: *pubKeyFromScalar(scalar),
	}
}

// pubKeyFromScalar returns a new public key instance.
func pubKeyFromScalar(scalar *big.Int) *PublicKey {
	return newG2Point().ScalarMult(g2Gen, scalar).ToAffine()
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
	return []byte{}
	//return newG1Point().ScalarMult(g1Gen, priv.Secret).Marshal()
}

/*
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
