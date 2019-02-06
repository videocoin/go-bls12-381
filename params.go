package bls

var (
	// X represents the BLS12-381 construction paramater
	X = hexToBigNumber("0xd201000000010000")

	// Q (381 bits) = (x - 1)2 ((x4 - x2 + 1) / 3) + x
	Q = hexToBigNumber("0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")

	// R (255 bits) = (x4 - x2 + 1)
	R = hexToBigNumber("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
)
