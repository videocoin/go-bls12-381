package bls12

import "math/big"

// twistPoint is a curve point in the elliptic curve's twist over an extension
// field Fq². T = z1². To obtain the full speed of pairings on Weierstrass
// curves it is useful to represent a point by (X1 : Y1 : Z1 : T1) with T1 = Z²
type twistPoint struct {
	x, y, z, t fq2
}

func newTwistPoint(x, y fq2) *twistPoint {
	return &twistPoint{
		x: x,
		y: y,
		z: fq2{c0: *new(fq).SetUint64(1)},
		t: fq2{c0: *new(fq).SetUint64(1)},
	}
}

// Set sets c to the value of a and returns c.
func (c *twistPoint) Set(a *twistPoint) *twistPoint {
	c.x, c.y, c.z, c.t = a.x, a.y, a.z, a.t
	return c
}

// Equal reports whether a is equal to b.
func (a *twistPoint) Equal(b *twistPoint) bool {
	return *a == *b
}

// IsInfinity reports whether the point is at infinity.
func (a *twistPoint) IsInfinity() bool {
	return a.z == fq2{}
}

// Add sets c to the sum a+b and returns c.
func (c *twistPoint) Add(a, b *twistPoint) *twistPoint {
	if a.IsInfinity() {
		return c.Set(b)
	}
	if b.IsInfinity() {
		return c.Set(a)
	}

	// faster than Add
	if a.Equal(b) {
		return c.Double(a)
	}

	// See https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	z1z1 := new(fq2).Sqr(&a.z)
	z2z2 := new(fq2).Sqr(&b.z)

	u1 := new(fq2).Mul(&a.x, z2z2)
	u2 := new(fq2).Mul(&b.x, z1z1)

	s1 := new(fq2).Mul(&a.y, &b.z)
	s1.Mul(s1, z2z2)
	s2 := new(fq2).Mul(&b.y, &a.z)
	s2.Mul(s2, z1z1)

	h := new(fq2).Sub(u2, u1)
	i := new(fq2).Add(h, h)
	i.Sqr(i)
	j := new(fq2).Mul(h, i)
	r := new(fq2).Sub(s2, s1)
	r.Add(r, r)
	v := new(fq2).Mul(u1, i)

	p, t0, t1 := new(twistPoint), new(fq2), new(fq2)
	t0.Add(v, v)
	t0.Add(t0, j)
	p.x.Sqr(r)
	p.x.Sub(&p.x, t0)

	t0.Add(s1, s1)
	t0.Mul(t0, j)
	t1.Sub(v, &p.x)
	t1.Mul(t1, r)
	p.y.Sub(t1, t0)

	p.z.Add(&a.z, &b.z)
	p.z.Sqr(&p.z)
	t0.Add(z1z1, z2z2)
	p.z.Sub(&p.z, t0)
	p.z.Mul(&p.z, h)

	return c.Set(p)
}

// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
func (c *twistPoint) Double(a *twistPoint) *twistPoint {
	d, e, f, g, h, i := new(fq2), new(fq2), new(fq2), new(fq2), new(fq2), new(fq2)
	d.Sqr(&a.x)
	e.Sqr(&a.y)
	f.Sqr(e)
	g.Mul(&a.x, e)
	g.Add(g, g)
	g.Add(g, g)
	h.Add(d, d)
	h.Add(h, d)
	i.Sqr(h)

	p, t0 := new(twistPoint), new(fq2)
	p.x.Add(g, g)
	p.x.Sub(i, &p.x)
	t0.Add(f, f)
	t0.Add(t0, t0)
	t0.Add(t0, t0)
	p.y.Sub(g, &p.x)
	p.y.Mul(&p.y, h)
	p.y.Sub(&p.y, t0)
	p.z.Mul(&a.y, &a.z)
	p.z.Add(&p.z, &p.z)

	return c.Set(p)
}

/*
// ScalarMult returns b*(Ax,Ay) where b is a number in big-endian form.
// See https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add.
func (c *twistPoint) ScalarMult(a *twistPoint, b *big.Int) *twistPoint {
	p := &twistPoint{}
	for i := b.BitLen() - 1; i >= 0; i-- {
		p.Double(p)
		if b.Bit(i) == 1 {
			p.Add(p, a)
		}
	}

	return c.Set(p)
}


// ScalarMult returns b*(Ax,Ay) where b is a number in big-endian form.
// ScalarMult implements the 4-GLS algorithm.
// See https://eprint.iacr.org/2011/608.pdf.
func (c *twistPoint) ScalarMult(a *twistPoint, b *big.Int) *twistPoint {
	// compute endomorphisms
	// Φ = ψϕψ^−1
	glvEnd := func(p *twistPoint) *twistPoint {
		p.x.ScalarMult(&p.x, &frobFq6C2[2].c0)
		return p
	}

	// Ψ = ψFrobpψ^−1
	// See https://github.com/cloudflare/bn256/blob/master/optate.go#L160
	glsEnd := func(p *twistPoint) *twistPoint {
		p.x.Conjugate(&p.x).ScalarMult(&p.x, &frobFq6C1[1].c0)
		p.y.Conjugate(&p.y).Mul(&p.y, xiToPMinus1Over2)
		return p
	}

	// precompute lookup table
	sum := [16]*twistPoint{}
	for i := 1; i < len(sum); i++ {
		sum[i] = new(twistPoint)
	}
	sum[1].Set(a)
	sum[2].Set(a)
	sum[4].Set(a)
	sum[8].Set(a)
	glvEnd(sum[2])
	glsEnd(sum[4])
	glsEnd(glvEnd(sum[8]))

	subScalars := glsLattice.Decompose(b)
	fmt.Println(subScalars)

	// make subscalars positive
	exp := 1
	for _, si := range subScalars {
		if si.Sign() == -1 {
			si.Neg(si)
			sum[exp].Inverse(sum[exp])
		}
		exp *= 2
	}
	fmt.Println(subScalars)

	// complete lookup table
	sum[3].Add(sum[1], sum[2])
	sum[5].Add(sum[4], sum[1])
	sum[6].Add(sum[4], sum[2])
	sum[7].Add(sum[6], sum[1])
	sum[9].Add(sum[8], sum[1])
	sum[10].Add(sum[8], sum[2])
	sum[11].Add(sum[10], sum[1])
	sum[12].Add(sum[8], sum[4])
	sum[13].Add(sum[12], sum[1])
	sum[14].Add(sum[12], sum[2])
	sum[15].Add(sum[14], sum[1])
	fmt.Println(sum)

	multiScalar := multiScalarRecoding(subScalars)
	r := new(twistPoint)
	for i := len(multiScalar) - 1; i >= 0; i-- {
		r.Double(r)
		if multiScalar[i] != 0 {
			r.Add(r, sum[multiScalar[i]])
		}
	}

	return c.Set(r)
}
*/

// ScalarMult returns b*(Ax,Ay) where b is a number in big-endian form.
// ScalarMult implements the 2-GLV algorithm.
// See Guide to Pairing-Based Cryptography - Algorithm 6.2.
// TODO mixed addition - precompute = affine?
func (c *twistPoint) ScalarMult(a *twistPoint, b *big.Int) *twistPoint {
	// precompute lookup table
	sum := [4]*twistPoint{
		nil,
		new(twistPoint).Set(a),
		new(twistPoint).Set(a),
		&twistPoint{}, // computed as soon as the final subscalars are known
	}
	sum[2].x.ScalarMult(&sum[2].x, &frobFq6C2[2].c0)

	subScalars := glvLattice.Decompose(b)

	// make subscalars positive
	exp := 1
	for _, si := range subScalars {
		if si.Sign() == -1 {
			si.Neg(si)
			sum[exp].Inverse(sum[exp])
		}
		exp *= 2
	}

	// complete lookup table
	sum[3].Add(sum[1], sum[2])

	multiScalar := multiScalarRecoding(subScalars)
	r := new(twistPoint)
	for i := len(multiScalar) - 1; i >= 0; i-- {
		r.Double(r)
		if multiScalar[i] != 0 {
			r.Add(r, sum[multiScalar[i]])
		}
	}

	return c.Set(r)
}

// ToAffine sets a to its affine value and returns a.
// See https://www.sciencedirect.com/topics/computer-science/affine-coordinate - Jacobian Projective Points
func (a *twistPoint) ToAffine() *twistPoint {
	if (a.z.c0 == *new(fq).SetUint64(1)) && (a.z.c1 == fq{}) {
		return a
	}

	zInv, zInvSqr, zInvCube := new(fq2), new(fq2), new(fq2)
	zInv.Inv(&a.z)
	zInvSqr.Sqr(zInv)
	zInvCube.Mul(zInvSqr, zInv)
	a.x.Mul(&a.x, zInvSqr)
	a.y.Mul(&a.y, zInvCube)
	a.z.SetOne()
	a.t.SetOne()

	return a
}

// Inverse sets c to -a and returns c.
func (c *twistPoint) Inverse(a *twistPoint) *twistPoint {
	c.x.Set(&a.x)
	c.y.Neg(&a.y)
	c.z.Set(&a.z)
	c.t.SetZero()
	return c
}
