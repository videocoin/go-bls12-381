package bls12

import (
	"fmt"
	"math/big"
)

var (
	bigHalfR      = new(big.Int).Rsh(r, 1)
	bigOne        = new(big.Int).SetUint64(1)
	v1            = bigFromBase10("57896044618658097711785492504343953926634992332820282019728792003956564819968")
	v1Rsh    uint = 255

	// See Guide to Pairing-Based Cryptography - Decompositions for the k = 12 BLS
	// family.

	glvLattice = &lattice{
		basis: [][]*big.Int{
			{bigFromBase10("228988810152649578064853576960394133503"), bigFromBase10("-1")},
			{bigFromBase10("1"), bigFromBase10("228988810152649578064853576960394133504")},
		},
		adj: []*big.Int{bigFromBase10("228988810152649578064853576960394133504"), bigFromBase10("1")},
	}

	glsLattice = &lattice{
		basis: [][]*big.Int{
			{bigFromBase10("15132376222941642752"), bigFromBase10("0"), bigFromBase10("1"), bigFromBase10("0")},
			{bigFromBase10("1"), bigFromBase10("1"), bigFromBase10("15132376222941642752"), bigFromBase10("0")},
			{bigFromBase10("0"), bigFromBase10("15132376222941642752"), bigFromBase10("-1"), bigFromBase10("1")},
			{bigFromBase10("0"), bigFromBase10("-1"), bigFromBase10("0"), bigFromBase10("-15132376222941642752")},
		},
		adj: []*big.Int{
			bigFromBase10("3465144826073652318776269530687742778270252468765361963008"),
			bigFromBase10("-228988810152649578064853576960394133503"),
			bigFromBase10("15132376222941642752"),
			bigFromBase10("1"),
		},
		li: []*big.Int{
			bigFromBase10("3825971794891275542975308155226117830230659266969465460776"),
			bigFromBase10("-252833509987073905889124965614887630617"),
			bigFromBase10("-16708116839162527449"),
			bigFromBase10("1"),
		},
	}
)

type lattice struct {
	basis [][]*big.Int
	adj   []*big.Int
	det   *big.Int
	li    []*big.Int // roundings
}

// Decompose implements Babai's rounding technique. Decompose decomposes n as
// m-dimensional expansions n≡n0+n1λ+···+n(m−1)λ^(m−1). An LLL-reduced lattice
// basis must be precomputed. See Pairing-Based Cryptography - Page 213/214 -
// The GLV method.
func (l *lattice) Decompose(n *big.Int) []*big.Int {
	t0, t1 := new(big.Int), new(big.Int)
	m := len(l.adj)

	// w = (n,0)
	// v ~ wB^-1
	v := make([]*big.Int, m)

	/*
		v[0] = new(big.Int)
		v[0].DivMod(t0.Mul(n, l.adj[0]), l.det, t1)
		round(v[0], t1)
		v[1] = new(big.Int)
		v[1].DivMod(n, l.det, t1)
		round(v[1], t1)
	*/

	for i := 0; i < m; i++ {
		v[i] = new(big.Int)
		v[i].DivMod(t0.Mul(n, l.adj[i]), r, t1)
		// v[i].DivMod(t0.Mul(n, l.li[i]), v1, t1)
		//v[i].Mul(n, l.li[i]).Rsh(v[i], v1Rsh)
		round(v[i], t1)
	}

	// u = w - vB
	u := make([]*big.Int, m)
	for i := 0; i < m; i++ {
		u[i] = new(big.Int)
		for j := 0; j < m; j++ {
			t0.Mul(v[j], l.basis[j][i])
			u[i].Add(u[i], t0)
		}
		u[i].Neg(u[i])
	}
	u[0].Add(u[0], n)

	return u
}

func round(n, m *big.Int) {
	if m.Cmp(bigHalfR) == 1 {
		if n.Sign() == -1 {
			fmt.Println("Entrou")
			n.Sub(n, bigOne)
		} else {
			n.Add(n, bigOne)
		}

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
