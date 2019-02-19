package bls12

import "math/big"

const decimalBase = 10

var (
	bigQ = bigFromBase10(q)
	big0 = new(big.Int).SetUint64(0)
	big1 = new(big.Int).SetUint64(1)
)

func bigFromBase10(str string) *big.Int {
	n, _ := new(big.Int).SetString(str, decimalBase)
	return n
}
