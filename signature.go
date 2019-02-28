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
			twistPoint: g2.element(index),
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
	//return blsSignature{curvePoint: new(curvePoint).mul(hashToCurveSubGroup(hash, g1), priv.Secret)}
	return []byte("")
}
