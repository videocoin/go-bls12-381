package bls12

// TODO package description
// This package operates, internally, on projective coordinates. For a given
// (x, y) position on the curve, the Jacobian coordinates are (x1, y1, z1)
// where x = x1/z1² and y = y1/z1³.
// See https://www.nayuki.io/page/elliptic-curve-point-addition-in-projective-coordinates.
// # BLS parameter, used to generate other parameters: x = -0xd201000000010000
// TODO replace fq2One
// TODO Benchmark sets and no sets, calls to lower levels that are not necessary

const (
	// uAbs is the absolute value of u.
	uAbs uint64 = 15132376222941642752

	// k64 is a pre-calculated quantity equal to k mod R where k=(r(r^−1 mod n)−1)/n
	k64 uint64 = 0x89f3fffcfffcfffd
)

var (
	// q is a prime number that specifies the number of elements of the finite field.
	q, _ = bigFromBase10("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787")

	// q64 is q as 64 bit words.
	q64 = [6]uint64{0xB9FEFFFFFFFFAAAB, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}

	// r2 is used as an optimization to enter and leave the Montgomery domain.
	// See http://home.deib.polimi.it/pelosi/lib/exe/fetch.php?media=teaching:montgomery.pdf page 12/17
	r2, _ = fqFromBase10("2708263910654730174793787626328176511836455197166317677006154293982164122222515399004018013397331347120527951271750")

	// Since the nonzero elements of GF(pn) form a finite group with respect to multiplication, apn−1 = 1 (for a ≠ 0), thus the inverse of a is a^pn−2.
	qMinusTwo, _ = fqFromBase10("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559785")

	// g1X is the x-coordinate of G1's generator in the Montgomery form
	g1X, _ = new(fq).SetString("3685416753713387016781088315183077757961620795782546409894578378688607592378376318836054947676345821548104185464507")

	// g1Y is the y-coordinate of G1's generator in the Montgomery form
	g1Y, _ = new(fq).SetString("1339506544944476473020471379941921221584933875938349620426543736416511423956333506472724655353366534992391756441569")

	// g2X0 is the c0 x-coordinate of G2's generator in the Montgomery form
	g2X0, _ = new(fq).SetString("352701069587466618187139116011060144890029952792775240219908644239793785735715026873347600343865175952761926303160")

	// g2X1 is the c1 x-coordinate of G2's generator in the Montgomery form
	g2X1, _ = new(fq).SetString("3059144344244213709971259814753781636986470325476647558659373206291635324768958432433509563104347017837885763365758")

	// g2Y0 is the c0 y-coordinate of G2's generator in the Montgomery form
	g2Y0, _ = new(fq).SetString("1985150602287291935568054521177171638300868978215655730859378665066344726373823718423869104263333984641494340347905")

	// g2Y1 is the c1 y-coordinate of G2's generator in the Montgomery form
	g2Y1, _ = new(fq).SetString("927553665492332455747201965776037880757740193453592970025027978793976877002675564980949289727957565575433344219582")

	// g1Cofactor is the cofactor by which to multiply points to map them to G1. (on to the r-torsion). h = (x - 1)2 / 3
	g1Cofactor, _ = bigFromBase10("76329603384216526031706109802092473003")

	// Fq2(u + 1)**(((p^power) - 1) / 6), power E [0, 11]
	frob12c1 = [12]*fq2{
		fq2One,
		&fq2{
			fq{0x7089552b319d465, 0xc6695f92b50a8313, 0x97e83cccd117228f, 0xa35baecab2dc29ee, 0x1ce393ea5daace4d, 0x8f2220fb0fb66eb},
			fq{0xb2f66aad4ce5d646, 0x5842a06bfc497cec, 0xcf4895d42599d394, 0xc11b9cba40a8e8d0, 0x2e3813cbe5a0de89, 0x110eefda88847faf},
		},
		&fq2{
			fq{0xecfb361b798dba3a, 0xc100ddb891865a2c, 0xec08ff1232bda8e, 0xd5c13cc6f1ca4721, 0x47222a47bf7b5c04, 0x110f184e51c5f59},
			fq{},
		},
		&fq2{
			fq{0x3e2f585da55c9ad1, 0x4294213d86c18183, 0x382844c88b623732, 0x92ad2afd19103e18, 0x1d794e4fac7cf0b9, 0xbd592fc7d825ec8},
			fq{0x7bcfa7a25aa30fda, 0xdc17dec12a927e7c, 0x2f088dd86b4ebef1, 0xd1ca2087da74d4a7, 0x2da2596696cebc1d, 0xe2b7eedbbfd87d2},
		},
		&fq2{
			fq{0x30f1361b798a64e8, 0xf3b8ddab7ece5a2a, 0x16a8ca3ac61577f7, 0xc26a2ff874fd029b, 0x3636b76660701c6e, 0x51ba4ab241b6160},
			fq{},
		},
		&fq2{
			fq{0x3726c30af242c66c, 0x7c2ac1aad1b6fe70, 0xa04007fbba4b14a2, 0xef517c3266341429, 0x95ba654ed2226b, 0x2e370eccc86f7dd},
			fq{0x82d83cf50dbce43f, 0xa2813e53df9d018f, 0xc6f0caa53c65e181, 0x7525cf528d50fe95, 0x4a85ed50f4798a6b, 0x171da0fd6cf8eebd},
		},
		&fq2{
			fq{0x43f5fffffffcaaae, 0x32b7fff2ed47fffd, 0x7e83a49a2e99d69, 0xeca8f3318332bb7a, 0xef148d1ea0f4c069, 0x40ab3263eff0206},
			fq{},
		},
		&fq2{
			fq{0xb2f66aad4ce5d646, 0x5842a06bfc497cec, 0xcf4895d42599d394, 0xc11b9cba40a8e8d0, 0x2e3813cbe5a0de89, 0x110eefda88847faf},
			fq{0x7089552b319d465, 0xc6695f92b50a8313, 0x97e83cccd117228f, 0xa35baecab2dc29ee, 0x1ce393ea5daace4d, 0x8f2220fb0fb66eb},
		},
		&fq2{
			fq{0xcd03c9e48671f071, 0x5dab22461fcda5d2, 0x587042afd3851b95, 0x8eb60ebe01bacb9e, 0x3f97d6e83d050d2, 0x18f0206554638741},
			fq{},
		},
		&fq2{
			fq{0x7bcfa7a25aa30fda, 0xdc17dec12a927e7c, 0x2f088dd86b4ebef1, 0xd1ca2087da74d4a7, 0x2da2596696cebc1d, 0xe2b7eedbbfd87d2},
			fq{0x3e2f585da55c9ad1, 0x4294213d86c18183, 0x382844c88b623732, 0x92ad2afd19103e18, 0x1d794e4fac7cf0b9, 0xbd592fc7d825ec8},
		},
		&fq2{
			fq{0x890dc9e4867545c3, 0x2af322533285a5d5, 0x50880866309b7e2c, 0xa20d1b8c7e881024, 0x14e4f04fe2db9068, 0x14e56d3f1564853a},
			fq{},
		},
		&fq2{
			fq{0x82d83cf50dbce43f, 0xa2813e53df9d018f, 0xc6f0caa53c65e181, 0x7525cf528d50fe95, 0x4a85ed50f4798a6b, 0x171da0fd6cf8eebd},
			fq{0x3726c30af242c66c, 0x7c2ac1aad1b6fe70, 0xa04007fbba4b14a2, 0xef517c3266341429, 0x95ba654ed2226b, 0x2e370eccc86f7dd},
		},
	}

	// Fq2(u + 1)**(((p^power) - 1) / 3), power E [0, 5]
	frob6c1 = [6]*fq2{
		// TODO fqOne?
		&fq2{
			fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
			fq{},
		},
		&fq2{
			fq{},
			fq{0xcd03c9e48671f071, 0x5dab22461fcda5d2, 0x587042afd3851b95, 0x8eb60ebe01bacb9e, 0x3f97d6e83d050d2, 0x18f0206554638741},
		},
		&fq2{
			fq{0x30f1361b798a64e8, 0xf3b8ddab7ece5a2a, 0x16a8ca3ac61577f7, 0xc26a2ff874fd029b, 0x3636b76660701c6e, 0x51ba4ab241b6160},
			fq{},
		},
		&fq2{
			fq{},
			fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
		},
		&fq2{
			fq{0xcd03c9e48671f071, 0x5dab22461fcda5d2, 0x587042afd3851b95, 0x8eb60ebe01bacb9e, 0x3f97d6e83d050d2, 0x18f0206554638741},
			fq{},
		},
		&fq2{
			fq{},
			fq{0x30f1361b798a64e8, 0xf3b8ddab7ece5a2a, 0x16a8ca3ac61577f7, 0xc26a2ff874fd029b, 0x3636b76660701c6e, 0x51ba4ab241b6160},
		},
	}

	// Fq2(u + 1)**(((2p^power) - 2) / 3), power E [0, 5]
	frob6c2 = [6]*fq2{
		// TODO fqOne?
		&fq2{
			fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
			fq{},
		},
		&fq2{
			fq{0x890dc9e4867545c3, 0x2af322533285a5d5, 0x50880866309b7e2c, 0xa20d1b8c7e881024, 0x14e4f04fe2db9068, 0x14e56d3f1564853a},
			fq{},
		},
		&fq2{
			fq{0xcd03c9e48671f071, 0x5dab22461fcda5d2, 0x587042afd3851b95, 0x8eb60ebe01bacb9e, 0x3f97d6e83d050d2, 0x18f0206554638741},
			fq{},
		},
		&fq2{
			fq{0x43f5fffffffcaaae, 0x32b7fff2ed47fffd, 0x7e83a49a2e99d69, 0xeca8f3318332bb7a, 0xef148d1ea0f4c069, 0x40ab3263eff0206},
			fq{},
		},
		&fq2{
			fq{0x30f1361b798a64e8, 0xf3b8ddab7ece5a2a, 0x16a8ca3ac61577f7, 0xc26a2ff874fd029b, 0x3636b76660701c6e, 0x51ba4ab241b6160},
			fq{},
		},
		&fq2{
			fq{0xecfb361b798dba3a, 0xc100ddb891865a2c, 0xec08ff1232bda8e, 0xd5c13cc6f1ca4721, 0x47222a47bf7b5c04, 0x110f184e51c5f59},
			fq{},
		},
	}

	// Fq(-1)**(((p^power) - 1) / 2), power E [0, 1]
	frob2c1 = [2]*fq{
		&fqOne,
		&fq{0x43f5fffffffcaaae, 0x32b7fff2ed47fffd, 0x7e83a49a2e99d69, 0xeca8f3318332bb7a, 0xef148d1ea0f4c069, 0x40ab3263eff0206},
	}
)
