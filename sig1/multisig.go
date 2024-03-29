// Package sig1 implements the BLS signature scheme with signatures on G1.
// See https://eprint.iacr.org/2018/483.pdf.
package sig1

import (
	"io"
	"math/big"

	bls12 "github.com/videocoin/go-bls12-381"
)

// PublicKey represents a BLS public key.
type PublicKey struct {
	bls12.G2Point
}

// Aggregate sets z to the sum x+y and returns z.
func (z *PublicKey) Aggregate(x, y *PublicKey) *PublicKey {
	z.Add(&x.G2Point, &y.G2Point)
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
	bls12.G1Point
}

// Aggregate sets z to the sum x+y and returns z.
func (z *Signature) Aggregate(x, y *Signature) {
	z.Add(&x.G1Point, &y.G1Point)
}

// GenerateKey generates a public and private key pair.
func GenerateKey(reader io.Reader) (*PrivateKey, error) {
	k, err := bls12.RandFieldElement(reader)
	if err != nil {
		return nil, err
	}

	priv := &PrivateKey{
		Secret:    k,
		PublicKey: PublicKey{*new(bls12.G2Point).ScalarBaseMult(k).ToAffine()},
	}

	return priv, nil
}

// Sign signs a hash using the private key, priv.
func Sign(priv *PrivateKey, hash []byte) *Signature {
	return &Signature{*new(bls12.G1Point).ScalarMult(new(bls12.G1Point).HashToPoint(hash), priv.Secret)}
}

// Verify verifies the signature of hash using the public key, pub. Its
// return value records whether the signature is valid.
func Verify(hash []byte, sig *Signature, pubKey *PublicKey) bool {
	return bls12.Pair(&sig.G1Point, bls12.G2Gen).Equal(bls12.Pair(new(bls12.G1Point).HashToPoint(hash), &pubKey.G2Point))
}

// VerifyAggregateCommon verifies that a signature is valid, for a collection
// of public keys and a common message. Its return value records whether the
// signature is valid.
func VerifyAggregateCommon(hash []byte, multiSig *Signature, pubKeys []*PublicKey) bool {
	return bls12.Pair(&multiSig.G1Point, bls12.G2Gen).Equal(bls12.Pair(new(bls12.G1Point).HashToPoint(hash), &AggregatePublicKeys(pubKeys).G2Point))
}

// VerifyAggregateDistinct verifies that a signature is valid, for a collection
// of public keys and distinct messages. Its return value records whether the
// signature is valid.
func VerifyAggregateDistinct(hashes [][]byte, multiSig *Signature, pubKeys []*PublicKey) bool {
	// TODO
	distinct := true

	if !distinct {
		return false
	}

	return VerifyAggregate(hashes, multiSig, pubKeys)
}

// VerifyAggregate verifies that a signature is valid, for a collection of
// public keys and messages. Its return value records whether the signature is
// valid. This method should only be used directly if the user proved knowledge
// or possesion of the corresponding secret key to prevent rogue public-key
// attacks. Message distinctness can be enforced by always prepending the public
// key to every message prior to signing. However, because now all messages are
// distinct, you cannot take advantage of VerifyAggregateCommon.
func VerifyAggregate(hashes [][]byte, multiSig *Signature, pubKeys []*PublicKey) bool {
	if len(hashes) != len(pubKeys) {
		return false
	}
	if len(hashes) == 0 {
		return false
	}

	t0 := new(bls12.G1Point).HashToPoint(hashes[0])
	pairing := bls12.Pair(t0, &pubKeys[0].G2Point)
	for i, pi := range pubKeys[1:] {
		t0.HashToPoint(hashes[i])
		pairing.Add(pairing, bls12.Pair(t0, &pi.G2Point))
	}

	return bls12.Pair(&multiSig.G1Point, bls12.G2Gen).Equal(pairing)
}

// AggregateSignatures aggregates multiple signatures into one signature.
func AggregateSignatures(sigs []*Signature) *Signature {
	sig := new(Signature)
	for _, si := range sigs {
		sig.Aggregate(sig, si)
	}
	return sig
}

// AggregatePublicKeys aggregates multiple public keys into one public key.
func AggregatePublicKeys(pubKeys []*PublicKey) *PublicKey {
	pub := new(PublicKey)
	for _, pi := range pubKeys {
		pub.Aggregate(pub, pi)
	}
	return pub
}
