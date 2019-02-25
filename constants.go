package bls12

const (
	// q is a prime number that specifies the number of elements of the finite field
	q = "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787"

	// k64 is a pre-calculated quantity equal to k mod R where k=(r(r^−1 mod n)−1)/n
	k64 = uint64(0x89f3fffcfffcfffd)
)

var (
	// r2 is used as an optimization to enter and leave the Montgomery domain
	// See http://home.deib.polimi.it/pelosi/lib/exe/fetch.php?media=teaching:montgomery.pdf page 12/17
	r2, _ = fqFromBig(bigFromBase10("2708263910654730174793787626328176511836455197166317677006154293982164122222515399004018013397331347120527951271750"))

	// q64 is q as 64 bit words
	q64 = [6]uint64{0xB9FEFFFFFFFFAAAB, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}
)
