package bls12

import (
	"fmt"
	"math/big"
)

var (
	// g2X0 is the c0 x-coordinate of G2's generator in the Montgomery form
	g2X0, _ = new(fq).SetString("352701069587466618187139116011060144890029952792775240219908644239793785735715026873347600343865175952761926303160", Montgomery)

	// g2X1 is the c1 x-coordinate of G2's generator in the Montgomery form
	g2X1, _ = new(fq).SetString("3059144344244213709971259814753781636986470325476647558659373206291635324768958432433509563104347017837885763365758", Montgomery)

	// g2Y0 is the c0 y-coordinate of G2's generator in the Montgomery form
	g2Y0, _ = new(fq).SetString("1985150602287291935568054521177171638300868978215655730859378665066344726373823718423869104263333984641494340347905", Montgomery)

	// g2Y1 is the c1 y-coordinate of G2's generator in the Montgomery form
	g2Y1, _ = new(fq).SetString("927553665492332455747201965776037880757740193453592970025027978793976877002675564980949289727957565575433344219582", Montgomery)

	g2Gen = &g2Point{*newTwistPoint(fq2{*g2X0, *g2X1}, fq2{*g2Y0, *g2Y1})}
)

type g2Point struct {
	p twistPoint
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (z *g2Point) ScalarBaseMult(scalar *big.Int) *g2Point {
	return z.ScalarMult(g2Gen, scalar)
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
func (z *g2Point) ScalarMult(x *g2Point, scalar *big.Int) *g2Point {
	z.p.ScalarMult(&x.p, scalar)
	return z
}

func (z *g2Point) Equal(x *g2Point) bool {
	return z.p.Equal(&x.p)
}

// Add returns the sum of (x1,y1) and (x2,y2).
func (z *g2Point) Add(x, y *g2Point) *g2Point {
	z.p.Add(&x.p, &y.p)
	return z
}

func (z *g2Point) Set(x *g2Point) *g2Point {
	if z != x {
		z.p.Set(&x.p)
	}
	return z
}

func (z *g2Point) ToAffine() *g2Point {
	z.p.ToAffine()
	return z
}

func (z *g2Point) String() string {
	return fmt.Sprintf("%v", z.p)
}
