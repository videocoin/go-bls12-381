package bls12

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func TestFeldman(t *testing.T) {
	threshold := uint64(2)
	numShares := uint64(4)
	verificationVec, shares, priv, err := CreateShares(rand.Reader, nil, threshold, numShares)
	if err != nil {
		t.Fatal(err)
	}

	// minimum required number of shares
	newPriv, err := PrivKeyFromShares(shares[:threshold])
	if err != nil {
		t.Fatal(err)
	}
	if newPriv.Secret.Cmp(priv.Secret) != 0 {
		t.Fatalf("expected: %s, got: %s\n", priv.Secret, newPriv.Secret)
	}

	// all shares
	newPriv, err = PrivKeyFromShares(shares)
	if err != nil {
		t.Fatal(err)
	}
	if newPriv.Secret.Cmp(priv.Secret) != 0 {
		t.Fatalf("expected: %s, got: %s\n", priv.Secret, newPriv.Secret)
	}

	// not enough shares
	newPriv, err = PrivKeyFromShares(shares[:threshold-1])
	if err != nil {
		t.Fatal(err)
	}
	if newPriv.Secret.Cmp(priv.Secret) == 0 {
		t.Fatalf("expected: %s, got: %s\n", priv.Secret, newPriv.Secret)
	}

	// generated shares must be valid
	for _, share := range shares {
		if err := VerifyShare(share, verificationVec); err != nil {
			t.Fatalf("Share is not valid: index %d, value %v\n", share.X, share.Y)
		}
	}

	// invalid share
	invalidShare := &Share{
		X: 3,
		Y: privKeyFromScalar(new(big.Int).SetUint64(6)),
	}
	if err := VerifyShare(invalidShare, verificationVec); err == nil {
		t.Fatal("pub keys must be different")
	}

}
