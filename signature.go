package bls12

import (
	"io"
	"math/big"
)

// The modified BLS multi-signature construction
// See https://crypto.stanford.edu/~dabo/pubs/papers/BLSmultisig.html
// Keep it simple for now - no signature aggregation

var unmarshalSignature = unmarshalG1Point

// PublicKey represents a BLS12-381 public key.
type PublicKey = g2Point

// PrivateKey represents a BLS12-381 private key.
type PrivateKey = big.Int

func PubKeyFromPrivKey(privKey *PrivateKey) *PublicKey {
	return G2.Element(privKey)
}

type blsSignature struct {
	*g1Point
}

// GenerateKey generates a public and private key pair.
func GenerateKey(reader io.Reader) (*PrivateKey, error) {
	/*
		elm, err := RandFieldElement(reader)
		if err != nil {
			return nil, err
		}

		return newPrivateKey(elm), nil
	*/
	return nil, nil
}

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
