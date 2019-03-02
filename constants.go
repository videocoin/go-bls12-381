package bls12

// TODO package description
// This package operates, internally, on projective coordinates. For a given
// (x, y) position on the curve, the Jacobian coordinates are (x1, y1, z1)
// where x = x1/z1² and y = y1/z1³.
// See https://www.nayuki.io/page/elliptic-curve-point-addition-in-projective-coordinates.
// # BLS parameter, used to generate other parameters: x = -0xd201000000010000

const (
	// k64 is a pre-calculated quantity equal to k mod R where k=(r(r^−1 mod n)−1)/n
	k64 uint64 = 0x89f3fffcfffcfffd

	// TODO
	orderBits = 381
	// TODO
	orderBytes = (orderBits + 7) / 8
)

var (
	// q is a prime number that specifies the number of elements of the finite field
	q = bigFromBase10("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787")

	// r2 is used as an optimization to enter and leave the Montgomery domain
	// See http://home.deib.polimi.it/pelosi/lib/exe/fetch.php?media=teaching:montgomery.pdf page 12/17
	r2, _ = fqFromBase10("2708263910654730174793787626328176511836455197166317677006154293982164122222515399004018013397331347120527951271750")

	// TODO - figure out how this value is calculated
	qm2, _ = fqMontgomeryFromBase10("0")

	// q64 is q as 64 bit words
	q64 = [6]uint64{0xB9FEFFFFFFFFAAAB, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}

	// g1X is the x-coordinate of G1's generator in the Montgomery form
	g1X, _ = fqMontgomeryFromBase10("3685416753713387016781088315183077757961620795782546409894578378688607592378376318836054947676345821548104185464507")

	// g1Y is the y-coordinate of G1's generator in the Montgomery form
	g1Y, _ = fqMontgomeryFromBase10("1339506544944476473020471379941921221584933875938349620426543736416511423956333506472724655353366534992391756441569")

	// g2X0 is the c0 x-coordinate of G2's generator in the Montgomery form
	g2X0, _ = fqMontgomeryFromBase10("352701069587466618187139116011060144890029952792775240219908644239793785735715026873347600343865175952761926303160")

	// g2X1 is the c1 x-coordinate of G2's generator in the Montgomery form
	g2X1, _ = fqMontgomeryFromBase10("3059144344244213709971259814753781636986470325476647558659373206291635324768958432433509563104347017837885763365758")

	// g2Y0 is the c0 y-coordinate of G2's generator in the Montgomery form
	g2Y0, _ = fqMontgomeryFromBase10("1985150602287291935568054521177171638300868978215655730859378665066344726373823718423869104263333984641494340347905")

	// g2Y1 is the c1 y-coordinate of G2's generator in the Montgomery form
	g2Y1, _ = fqMontgomeryFromBase10("927553665492332455747201965776037880757740193453592970025027978793976877002675564980949289727957565575433344219582")

	// g1Cofactor is the cofactor by which to multiply points to map them to G1. (on to the r-torsion). h = (x - 1)2 / 3
	g1Cofactor = bigFromBase10("76329603384216526031706109802092473003")
)
