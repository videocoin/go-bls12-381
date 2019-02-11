package bls12

import "math/big"

const decimalBase = 10

func bigFromBase10(str string) *big.Int {
	n, _ := new(big.Int).SetString(str, decimalBase)
	return n
}
