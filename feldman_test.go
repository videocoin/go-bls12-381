package bls12

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func TestFeldman(t *testing.T) {
	threshold := uint64(6)
	numShares := uint64(20)
	verificationVec, shares, privKey, err := CreateShares(rand.Reader, threshold, numShares)
	if err != nil {
		t.Fatal(err)
	}

	// minimum required number of shares
	newPrivKey, err := PrivKeyFromShares(shares[:threshold])
	if err != nil {
		t.Fatal(err)
	}
	if newPrivKey.Secret.Cmp(privKey.Secret) != 0 {
		t.Errorf("expected: %s, got: %s\n", privKey.Secret, newPrivKey.Secret)
	}

	// all shares
	newPrivKey, err = PrivKeyFromShares(shares)
	if err != nil {
		t.Fatal(err)
	}
	if newPrivKey.Secret.Cmp(privKey.Secret) != 0 {
		t.Errorf("expected: %s, got: %s\n", privKey.Secret, newPrivKey.Secret)
	}

	// not enough shares
	newPrivKey, err = PrivKeyFromShares(shares[:threshold-1])
	if err != nil {
		t.Fatal(err)
	}
	if newPrivKey.Secret.Cmp(privKey.Secret) == 0 {
		t.Errorf("expected: %s, got: %s\n", privKey.Secret, newPrivKey.Secret)
	}

	// generated shares must be valid
	for _, share := range shares {
		if err := VerifyShare(share, verificationVec); err != nil {
			t.Fatal(err)
		}
	}

	invalidShare := &Share{
		X: 3,
		Y: privKeyFromScalar(new(big.Int).SetUint64(6)),
	}
	if err := VerifyShare(invalidShare, verificationVec); err == nil {
		t.Fatal("pub keys must be different")
	}
}

func TestCreateShares(t *testing.T) {
	// TODO
}

func TestPrivKeyFromShares(t *testing.T) {
	// TODO
}

func TestVerifyShare(t *testing.T) {
	// TODO
}
