package bls12

const (
	// q is a prime number that specifies the number of elements of the finite field (order)
	q = "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787"

	/*
		_WORD64 = 64        // word size in bits
		_WORD63 = _WORD64 - 1  // word size in bits - 1
		_WORD32 = _WORD64 / 2  // half word size in bits
		_BASE32 = 1 << _W64 // half digit base
		_SIZE32 = 32
		_MASK32 = _B32 - 1  // half digit mask
	*/
)

var (
	// r2 is used as an optimization to enter and leave the Montgomery domain.
	// See http://home.deib.polimi.it/pelosi/lib/exe/fetch.php?media=teaching:montgomery.pdf page 12/17
	r2, _ = fqFromBig(bigFromBase10("2708263910654730174793787626328176511836455197166317677006154293982164122222515399004018013397331347120527951271750"))

	// q64 is q as 64 bit words
	q64 = [6]uint64{0xB9FEFFFFFFFFAAAB, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}
)
