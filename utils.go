package bls

import "math/big"

const (
	hexBase = 16
)

func hexToBigNumber(hex string) *big.Int {
	bigNumber, _ := new(big.Int).SetString(hex, hexBase)
	return bigNumber
}
