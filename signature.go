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

type blsSignature struct {
	*curvePoint
}

func newPrivateKey(index *big.Int) *PrivateKey {
	return &PrivateKey{
		Secret: index,
		PublicKey: PublicKey{
			twistPoint: g2.getElement(index),
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
