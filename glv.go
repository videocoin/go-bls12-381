package bls12

import "math/big"

var (
	bigHalfR = new(big.Int).Rsh(r, 1)
	bigOne   = new(big.Int).SetUint64(1)

	// See Guide to Pairing-Based Cryptography - Decompositions for the k = 12 BLS
	// family.
	g1Lattice = &lattice{
		basis: [][]*big.Int{
			{bigFromBase10("228988810152649578064853576960394133503"), bigFromBase10("-1")},
			{bigFromBase10("1"), bigFromBase10("228988810152649578064853576960394133504")},
		},
		adj: []*big.Int{bigFromBase10("228988810152649578064853576960394133504"), bigFromBase10("-1")}, // TODO optmization -1?
		det: bigFromBase10("52435875175126190479447740508185965837690552500527637822603658699938581184513"),
	}
)

type lattice struct {
	basis [][]*big.Int
	adj   []*big.Int
	det   *big.Int
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
	for i := 0; i < m; i++ {
		v[i] = new(big.Int)
		v[i].DivMod(t0.Mul(n, l.adj[i]), l.det, t1)
		// round
		if t1.Cmp(bigHalfR) == 1 {
			v[i].Add(v[i], bigOne)
		}
	}

	// u = w - vB
	u := make([]*big.Int, m)
	for i := 0; i < m; i++ {
		u[i] = new(big.Int)
		for j := 0; j < m; j++ {
			t0.Mul(v[j], l.basis[i][j])
			u[i].Add(u[i], t0)
		}
		u[i] = u[i].Neg(u[i])
	}
	u[0].Add(u[0], n)

	return u
}

func multiScalarRecoding(scalars []*big.Int) []uint8 {
	// find the max bit length
	max := new(big.Int)
	for _, si := range scalars {
		if si.Cmp(max) == 1 {
			max.Set(si)
		}
	}

	// merge lattice vectors according to the pre computed sum
	multi := make([]uint8, max.BitLen())
	for i, si := range scalars {
		for j := 0; j < si.BitLen(); j++ {
			multi[j] += uint8(si.Bit(j)) << uint(i)
		}
	}

	return multi
}
