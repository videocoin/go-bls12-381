// Package bls12 implements bls12-381 pairing-friendly elliptic curve
// construction. This package operates, internally, on projective coordinates.
// For a given (x, y) position on the curve, the Jacobian coordinates are
// (x1, y1, z1) where x = x1/z1² and y = y1/z1³.
package bls12

import "math/big"

var (
	// q is a prime number that specifies the number of elements of the finite field.
	q, _ = bigFromBase10("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787")

	// q64 is q as 64 bit words.
	q64 = [6]uint64{0xB9FEFFFFFFFFAAAB, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}

	// qMinusTwo is the value by which to exponentiate q-order field elements to
	// calculate their inverse.
	qMinusTwo = &fq{0xB9FEFFFFFFFFAAA9, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}

	// r2Q is the value by which to multiply q-order field elements to map them to the Montgomery domain.
	qR2 = &fq{0xf4df1f341c341746, 0x0a76e6a609d104f1, 0x8de5476c4c95b6d5, 0x67eb88a9939d83c0, 0x9a793e85b519952d, 0x11988fe592cae3aa}

	// TODO review desc.
	// qK64 is a pre-calculated quantity equal to k mod R where k=(r(r^−1 mod n)−1)/n.
	qK64 uint64 = 0x89f3fffcfffcfffd

	// r is the order of the groups.
	r, _ = bigFromBase10("52435875175126190479447740508185965837690552500527637822603658699938581184513")

	// r64 is r as 64 bit words.
	r64 = [4]uint64{0xFFFFFFFF00000001, 0x53BDA402FFFE5BFE, 0x3339D80809A1D805, 0x73EDA753299D7D48}

	// TODO review desc.
	// rK64 is a pre-calculated quantity equal to k mod R where k=(r(r^−1 mod n)−1)/n.
	rK64 uint64 = 0xfffffffeffffffff

	// rMinusTwo is the value by which to exponentiate r-order field elements to
	// calculate their inverse.
	rMinusTwo = &fr{0xFFFFFFFEFFFFFFFF, 0x53BDA402FFFE5BFE, 0x3339D80809A1D805, 0x73EDA753299D7D48}

	// r2R is the value by which to multiply r-order field elements to map them to the Montgomery domain.
	rR2 = &fr{0xC999E990F3F29C6D, 0x2B6CEDCB87925C23, 0x05D314967254398F, 0x748D9D99F59FF11}

	// g1X is the x-coordinate of G1's generator.
	g1X = &fq{6679831729115696150, 8653662730902241269, 1535610680227111361, 17342916647841752903, 17135755455211762752, 1297449291367578485}

	// g1Y is the y-coordinate of G1's generator.
	g1Y = &fq{13451288730302620273, 10097742279870053774, 15949884091978425806, 5885175747529691540, 1016841820992199104, 845620083434234474}

	// g10 is used during the convertion of bytes to a point in g1.
	g10 = []byte("G1_0")

	// g10 is used during the convertion of bytes to a point in g1.
	g11 = []byte("G1_1")

	// g2X0 is the c0 x-coordinate of G2's generator.
	g2X0 = &fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735}

	// g2X1 is the c1 x-coordinate of G2's generator.
	g2X1 = &fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155}

	// g2Y0 is the c0 y-coordinate of G2's generator.
	g2Y0 = &fq{5508758831087832138, 6448303779119275098, 16710190169160573786, 13542242618704742751, 563980702369916322, 37152010398653157}

	// g2Y1 is the c1 y-coordinate of G2's generator.
	g2Y1 = &fq{12520284671833321565, 1777275927576994268, 9704602344324656032, 8739618045342622522, 16651875250601773805, 804950956836789234}

	// g1Cofactor is the cofactor by which to multiply points to map them to G1. (on to the r-torsion). h = (x - 1)2 / 3
	g1Cofactor, _ = bigFromBase10("76329603384216526031706109802092473003")

	// frobFq2C1 contains the value by which to multiply c1 to calculate the frobenius for a certain power.
	frobFq2C1 = [2]*fq{
		&fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
		&fq{0x43f5fffffffcaaae, 0x32b7fff2ed47fffd, 0x7e83a49a2e99d69, 0xeca8f3318332bb7a, 0xef148d1ea0f4c069, 0x40ab3263eff0206},
	}

	// frobFq6C1 contains the value by which to multiply c1 to calculate the frobenius for a certain power.
	frobFq6C1 = [6]*fq2{
		&fq2{
			c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
		},
		&fq2{
			c1: fq{0xcd03c9e48671f071, 0x5dab22461fcda5d2, 0x587042afd3851b95, 0x8eb60ebe01bacb9e, 0x3f97d6e83d050d2, 0x18f0206554638741},
		},
		&fq2{
			c0: fq{0x30f1361b798a64e8, 0xf3b8ddab7ece5a2a, 0x16a8ca3ac61577f7, 0xc26a2ff874fd029b, 0x3636b76660701c6e, 0x51ba4ab241b6160},
		},
		&fq2{
			c1: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
		},
		&fq2{
			c0: fq{0xcd03c9e48671f071, 0x5dab22461fcda5d2, 0x587042afd3851b95, 0x8eb60ebe01bacb9e, 0x3f97d6e83d050d2, 0x18f0206554638741},
		},
		&fq2{
			c1: fq{0x30f1361b798a64e8, 0xf3b8ddab7ece5a2a, 0x16a8ca3ac61577f7, 0xc26a2ff874fd029b, 0x3636b76660701c6e, 0x51ba4ab241b6160},
		},
	}

	// frobFq6C2 contains the value by which to multiply c2 to calculate the frobenius for a certain power.
	frobFq6C2 = [6]*fq2{
		&fq2{
			c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
		},
		&fq2{
			c0: fq{0x890dc9e4867545c3, 0x2af322533285a5d5, 0x50880866309b7e2c, 0xa20d1b8c7e881024, 0x14e4f04fe2db9068, 0x14e56d3f1564853a},
		},
		&fq2{
			c0: fq{0xcd03c9e48671f071, 0x5dab22461fcda5d2, 0x587042afd3851b95, 0x8eb60ebe01bacb9e, 0x3f97d6e83d050d2, 0x18f0206554638741},
		},
		&fq2{
			c0: fq{0x43f5fffffffcaaae, 0x32b7fff2ed47fffd, 0x7e83a49a2e99d69, 0xeca8f3318332bb7a, 0xef148d1ea0f4c069, 0x40ab3263eff0206},
		},
		&fq2{
			c0: fq{0x30f1361b798a64e8, 0xf3b8ddab7ece5a2a, 0x16a8ca3ac61577f7, 0xc26a2ff874fd029b, 0x3636b76660701c6e, 0x51ba4ab241b6160},
		},
		&fq2{
			c0: fq{0xecfb361b798dba3a, 0xc100ddb891865a2c, 0xec08ff1232bda8e, 0xd5c13cc6f1ca4721, 0x47222a47bf7b5c04, 0x110f184e51c5f59},
		},
	}

	// frobFq12C1 contains the value by which to multiply c1 to calculate the frobenius for a certain power.
	frobFq12C1 = [12]*fq2{
		&fq2{
			c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
		},
		&fq2{
			fq{0x7089552b319d465, 0xc6695f92b50a8313, 0x97e83cccd117228f, 0xa35baecab2dc29ee, 0x1ce393ea5daace4d, 0x8f2220fb0fb66eb},
			fq{0xb2f66aad4ce5d646, 0x5842a06bfc497cec, 0xcf4895d42599d394, 0xc11b9cba40a8e8d0, 0x2e3813cbe5a0de89, 0x110eefda88847faf},
		},
		&fq2{
			c0: fq{0xecfb361b798dba3a, 0xc100ddb891865a2c, 0xec08ff1232bda8e, 0xd5c13cc6f1ca4721, 0x47222a47bf7b5c04, 0x110f184e51c5f59},
		},
		&fq2{
			fq{0x3e2f585da55c9ad1, 0x4294213d86c18183, 0x382844c88b623732, 0x92ad2afd19103e18, 0x1d794e4fac7cf0b9, 0xbd592fc7d825ec8},
			fq{0x7bcfa7a25aa30fda, 0xdc17dec12a927e7c, 0x2f088dd86b4ebef1, 0xd1ca2087da74d4a7, 0x2da2596696cebc1d, 0xe2b7eedbbfd87d2},
		},
		&fq2{
			c0: fq{0x30f1361b798a64e8, 0xf3b8ddab7ece5a2a, 0x16a8ca3ac61577f7, 0xc26a2ff874fd029b, 0x3636b76660701c6e, 0x51ba4ab241b6160},
		},
		&fq2{
			fq{0x3726c30af242c66c, 0x7c2ac1aad1b6fe70, 0xa04007fbba4b14a2, 0xef517c3266341429, 0x95ba654ed2226b, 0x2e370eccc86f7dd},
			fq{0x82d83cf50dbce43f, 0xa2813e53df9d018f, 0xc6f0caa53c65e181, 0x7525cf528d50fe95, 0x4a85ed50f4798a6b, 0x171da0fd6cf8eebd},
		},
		&fq2{
			c0: fq{0x43f5fffffffcaaae, 0x32b7fff2ed47fffd, 0x7e83a49a2e99d69, 0xeca8f3318332bb7a, 0xef148d1ea0f4c069, 0x40ab3263eff0206},
		},
		&fq2{
			fq{0xb2f66aad4ce5d646, 0x5842a06bfc497cec, 0xcf4895d42599d394, 0xc11b9cba40a8e8d0, 0x2e3813cbe5a0de89, 0x110eefda88847faf},
			fq{0x7089552b319d465, 0xc6695f92b50a8313, 0x97e83cccd117228f, 0xa35baecab2dc29ee, 0x1ce393ea5daace4d, 0x8f2220fb0fb66eb},
		},
		&fq2{
			c0: fq{0xcd03c9e48671f071, 0x5dab22461fcda5d2, 0x587042afd3851b95, 0x8eb60ebe01bacb9e, 0x3f97d6e83d050d2, 0x18f0206554638741},
		},
		&fq2{
			fq{0x7bcfa7a25aa30fda, 0xdc17dec12a927e7c, 0x2f088dd86b4ebef1, 0xd1ca2087da74d4a7, 0x2da2596696cebc1d, 0xe2b7eedbbfd87d2},
			fq{0x3e2f585da55c9ad1, 0x4294213d86c18183, 0x382844c88b623732, 0x92ad2afd19103e18, 0x1d794e4fac7cf0b9, 0xbd592fc7d825ec8},
		},
		&fq2{
			c0: fq{0x890dc9e4867545c3, 0x2af322533285a5d5, 0x50880866309b7e2c, 0xa20d1b8c7e881024, 0x14e4f04fe2db9068, 0x14e56d3f1564853a},
		},
		&fq2{
			fq{0x82d83cf50dbce43f, 0xa2813e53df9d018f, 0xc6f0caa53c65e181, 0x7525cf528d50fe95, 0x4a85ed50f4798a6b, 0x171da0fd6cf8eebd},
			fq{0x3726c30af242c66c, 0x7c2ac1aad1b6fe70, 0xa04007fbba4b14a2, 0xef517c3266341429, 0x95ba654ed2226b, 0x2e370eccc86f7dd},
		},
	}
)

func bigFromBase10(s string) (*big.Int, bool) {
	return new(big.Int).SetString(s, decimalBase)
}
