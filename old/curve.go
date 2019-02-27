package bls12

import (
	"bytes"
	"math/big"
)

var curveB = newFq(bigFromBase10("4"))

// curvePoint is an elliptic curve(y²=x³+3) point over the finite field Fq.
type curvePoint struct {
	x, y, z fq
}

func newCurvePoint(x, y fq) *curvePoint {
	return &curvePoint{
		x: x,
		y: y,
	}
}

func (cp *curvePoint) isInfinity() bool {
	return cp.z.isZero()
}

func (cp *curvePoint) set(p *curvePoint) {
	cp.x = p.x
	cp.y = p.y
	cp.z = p.z
}

func (cp *curvePoint) mul(p *curvePoint, scalar *big.Int) *curvePoint {

	q := new(curvePoint)
	for i := scalar.BitLen(); i > 0; i-- {
		q.double(q)
		if scalar.Bit(i) != 0 {
			//q.add(q, p)
		}
	}
	cp.set(q)

	return cp
}

func (cp *curvePoint) Marshal() ([]byte, error) {
	// See https://github.com/zkcrypto/pairing/tree/master/src/bls12_381 - "Serialization"
	buffer := new(bytes.Buffer)

	x, err := cp.x.Marshal()
	if err != nil {
		return nil, err
	}
	buffer.Write(x)

	y, err := cp.y.Marshal()
	if err != nil {
		return nil, err
	}
	buffer.Write(y)

	return buffer.Bytes(), nil
}
