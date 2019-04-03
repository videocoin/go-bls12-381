package bls12

import (
	"math/big"
	"testing"
)

func TestSign(t *testing.T) {}

func TestAggregate(t *testing.T) {}

func TestVerify(t *testing.T) {}

func TestPubKeyFromPrivKey(t *testing.T) {
	testCases := []struct {
		privKey *PrivateKey
		pubKey  *PublicKey
	}{
		{
			privKey: new(big.Int).SetUint64(1),
			pubKey:  g2Generator,
		},
	}
	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			result := PubKeyFromPrivKey(testCase.privKey)
			if !result.Equal(testCase.pubKey) {
				t.Errorf("expected %v, got %v\n", testCase.pubKey, result)
			}
		})
	}
}
