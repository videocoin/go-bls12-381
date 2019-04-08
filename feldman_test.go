package bls12

import (
	"crypto/rand"
	"testing"
)

func TestFeldman(t *testing.T) {
	threshold := uint64(2)
	numShares := uint64(4)
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

	if err := VerifyShare(shares[0], verificationVec); err != nil {
		t.Fatal(err)
	}
	/*
		for _, share := range shares {

		}
	*/
}
