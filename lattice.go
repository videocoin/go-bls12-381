package bls12

import (
	"math/big"
)

const roundingRsh = 255

var (
	roundingInt = bigFromBase10("57896044618658097711785492504343953926634992332820282019728792003956564819968")

	// See Guide to Pairing-Based Cryptography Chapter 6.3.2
	glvLattice = &lattice{
		base: [][]*big.Int{
			{bigFromBase10("228988810152649578064853576960394133503"), bigFromBase10("-1")},
			{bigFromBase10("1"), bigFromBase10("228988810152649578064853576960394133504")},
		},
		adj: []*big.Int{bigFromBase10("228988810152649578064853576960394133504"), bigFromBase10("1")},
		det: r,
	}

	// See https://eprint.iacr.org/2013/458.pdf - 4-GLS on G2 for BN and BLS curves with k = 12
	glsLattice = &lattice{
		base: [][]*big.Int{
			{bigFromBase10("-15132376222941642752"), bigFromBase10("1"), bigFromBase10("0"), bigFromBase10("0")},
			{bigFromBase10("0"), bigFromBase10("-15132376222941642752"), bigFromBase10("1"), bigFromBase10("0")},
			{bigFromBase10("0"), bigFromBase10("0"), bigFromBase10("-15132376222941642752"), bigFromBase10("1")},
			{bigFromBase10("1"), bigFromBase10("0"), bigFromBase10("-1"), bigFromBase10("15132376222941642752")},
		},
		adj: []*big.Int{
			bigFromBase10("-3465144826073652318776269530687742778285384844988303605760"),
			bigFromBase10("-228988810152649578064853576960394133505"),
			bigFromBase10("-15132376222941642752"),
			bigFromBase10("-1"),
		},
		det: r,
	}
)

type lattice struct {
	base [][]*big.Int
	adj  []*big.Int
	det  *big.Int
}

func (l *lattice) Decompose(n *big.Int) []*big.Int {
	t0 := new(big.Int)
	m := len(l.adj)

	// w = (n,0)
	// v ~ wB^-1
	v := make([]*big.Int, m)
	for i := 0; i < m; i++ {
		v[i] = new(big.Int).Mul(n, l.adj[i])
		round(v[i], l.det)
	}
	//v[3] = new(big.Int)

	// u = w - vB
	u := make([]*big.Int, m)
	for i := 0; i < m; i++ {
		u[i] = new(big.Int)
		for j := 0; j < m; j++ {
			t0.Mul(v[j], l.base[j][i])
			u[i].Add(u[i], t0)
		}
		u[i].Neg(u[i])
	}
	u[0].Add(u[0], n)

	return u
}

// round sets num to num/denom rounded to the nearest integer.
func round(num, denom *big.Int) {
	r := new(big.Int)
	num.DivMod(num, denom, r)

	if r.Cmp(bigHalfR) == 1 {
		num.Add(num, big.NewInt(1))
	}
}

func multiScalarRecoding(scalars []*big.Int) []uint8 {
	max := new(big.Int)
	for _, si := range scalars {
		if si.Cmp(max) == 1 {
			max.Set(si)
		}
	}

	multi := make([]uint8, max.BitLen())
	for i, si := range scalars {
		for j := 0; j < si.BitLen(); j++ {
			multi[j] += uint8(si.Bit(j)) << uint(i)
		}
	}

	return multi
}
