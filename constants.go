package bls12

const (
	// q is a prime number that specifies the number of elements of the finite field (order)
	q  = "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787"
	r2 = "2708263910654730174793787626328176511836455197166317677006154293982164122222515399004018013397331347120527951271750"
	nq = ""
)

var (
	qU64 = [6]uint64{0xB9FEFFFFFFFFAAAB, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}
)

const (
	_W       = 64       // word size in bits
	_WMinus1 = _W - 1   // word size in bits - 1
	_W2      = _W / 2   // half word size in bits
	_B2      = 1 << _W2 // half digit base
	_M2      = _B2 - 1  // half digit mask
)
