package bls12

import (
	"math/big"
)

var (
	// Values taken from the execution of https://eprint.iacr.org/2019/403.pdf - A The isogeny maps.
	iso3XNum = []*fq2{
		&fq2{
			// 889424345604814976315064405719089812568196182208668418962679585805340366775741747653930584250892369786198727235542
			c0: fq{},
			// 889424345604814976315064405719089812568196182208668418962679585805340366775741747653930584250892369786198727235542
			c1: fq{},
		},
		&fq2{
			// 0
			c0: fq{},
			// 2668273036814444928945193217157269437704588546626005256888038757416021100327225242961791752752677109358596181706522
			c1: fq{},
		},
		&fq2{
			// 2668273036814444928945193217157269437704588546626005256888038757416021100327225242961791752752677109358596181706526
			c0: fq{},
			// 1334136518407222464472596608578634718852294273313002628444019378708010550163612621480895876376338554679298090853261
			c1: fq{},
		},
		&fq2{
			// 3557697382419259905260257622876359250272784728834673675850718343221361467102966990615722337003569479144794908942033
			c0: fq{},
			// 0
			c1: fq{},
		},
	}

	iso3XDen = []*fq2{
		&fq2{
			c0: fq{},
			// 4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559715
			c1: fq{},
		},
		&fq2{
			// 12
			c0: fq{},
			// 4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559775
			c1: fq{},
		},
		&fq2{
			// 1
			c0: fq{},
			c1: fq{},
		},
	}

	iso3YNum = []*fq2{
		&fq2{
			// 3261222600550988246488569487636662646083386001431784202863158481286248011511053074731078808919938689216061999863558
			c0: fq{},
			// 3261222600550988246488569487636662646083386001431784202863158481286248011511053074731078808919938689216061999863558
			c1: fq{},
		},
		&fq2{
			c0: fq{},
			// 889424345604814976315064405719089812568196182208668418962679585805340366775741747653930584250892369786198727235518
			c1: fq{},
		},
		&fq2{
			// 2668273036814444928945193217157269437704588546626005256888038757416021100327225242961791752752677109358596181706524
			c0: fq{},
			// 1334136518407222464472596608578634718852294273313002628444019378708010550163612621480895876376338554679298090853263
			c1: fq{},
		},
		&fq2{
			// 2816510427748580758331037284777117739799287910327449993381818688383577828123182200904113516794492504322962636245776
			c0: fq{},
			c1: fq{},
		},
	}

	iso3YDen = []*fq2{
		&fq2{
			// 4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559355
			c0: fq{},
			// 4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559355
			c1: fq{},
		},
		&fq2{
			c0: fq{},
			// 4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559571
			c1: fq{},
		},
		&fq2{
			// 18
			c0: fq{},
			// 4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559769
			c1: fq{},
		},
		&fq2{
			// 1
			c0: fq{},
			c1: fq{},
		},
	}

	// Values taken from the execution of https://eprint.iacr.org/2019/403.pdf - A The isogeny maps.
	iso3K = [][]*fq2{iso3XNum, iso3XDen, iso3YNum, iso3YDen}
)

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

func (a *twistPoint) iso3(b *twistPoint) *twistPoint {
	term := new(fq2)
	mul := new(fq2)
	var sum [4]fq2
	for i, ki := range iso3K {
		sum[i].Set(ki[0])
		mul.SetOne()
		for _, kij := range ki[1:] {
			mul.Mul(mul, &b.x)
			term.Mul(kij, mul)
			sum[i].Add(&sum[i], term)
		}
	}

	sum[1].Inv(&sum[1])
	a.x.Mul(&sum[0], &sum[1])
	sum[3].Inv(&sum[3])
	term.Mul(&sum[2], &sum[3])
	a.y.Mul(&b.y, term)

	return a
}

// TODO
// SetBytes sets c to the twist point that results from the given slice of bytes
// and returns c. The point is not guaranteed to be in a particular subgroup.
// See https://github.com/ethereum/eth2.0-specs/blob/dev/specs/bls_signature.md.
/*
func (c *curvePoint) SetBytes(buf []byte, ref0 []byte, ref1 []byte) *curvePoint {
	h := blake2b.Sum256(buf)
	sum := blake2b.Sum512(append(h[:], ref0...))
	t0 := new(big.Int)
	g10, _ := new(fq).SetInt(t0.Mod(t0.SetBytes(sum[:]), r))
	sum = blake2b.Sum512(append(h[:], ref1...))
	g11, _ := new(fq).SetInt(t0.Mod(t0.SetBytes(sum[:]), r))

	return c.Add(new(curvePoint).SWEncode(g10), new(curvePoint).SWEncode(g11))
}
*/
