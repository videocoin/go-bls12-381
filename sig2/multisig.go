// Package sig2 implements the BLS multi signature scheme with signatures on G2.
// See https://crypto.stanford.edu/~dabo/pubs/papers/BLSmultisig.html.
package sig2

import (
	"io"
	"math/big"

	bls12 "github.com/videocoin/go-bls12-381"
)

// PublicKey represents a BLS public key.
type PublicKey struct {
	bls12.G1Point
}

// Aggregate sets z to the sum x+y and returns z.
func (z *PublicKey) Aggregate(x, y *PublicKey) *PublicKey {
	z.Add(&x.G1Point, &y.G1Point)
	return z
}

// PrivateKey represents a BLS private key.
type PrivateKey struct {
	PublicKey
	Secret *big.Int
}

// Public returns the public key corresponding to priv.
func (priv *PrivateKey) Public() PublicKey {
	return priv.PublicKey
}

// Signature represents a BLS signature.
type Signature struct {
	bls12.G2Point
}

// Aggregate sets z to the sum x+y and returns z.
func (z *Signature) Aggregate(x, y *Signature) {
	z.Add(&x.G2Point, &y.G2Point)
}

func privKeyFromScalar(k *big.Int) *PrivateKey {
	priv := new(PrivateKey)
	priv.Secret = k
	priv.PublicKey = PublicKey{*new(bls12.G1Point).ScalarBaseMult(k).ToAffine()}
	return priv
}

// GenerateKey generates a public and private key pair.
func GenerateKey(reader io.Reader) (*PrivateKey, error) {
	k, err := bls12.RandFieldElement(reader)
	if err != nil {
		return nil, err
	}

	return privKeyFromScalar(k), nil
}

// Sign signs a hash using the private key, priv.
func Sign(priv *PrivateKey, hash []byte) []byte {
	return new(bls12.G2Point).ScalarMult(new(bls12.G2Point).HashToPoint(hash), priv.Secret).Marshal()
}
