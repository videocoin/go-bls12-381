package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

const fqLen = 6

var q64 = [6]uint64{0xB9FEFFFFFFFFAAAB, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}

func fqLoad(src Mem) [fqLen]Register {
	regs := [fqLen]Register{GP64(), GP64(), GP64(), GP64(), GP64(), GP64()}
	for i, ri := range regs {
		MOVQ(src.Offset(8*i), ri)
	}
	return regs
}

func fqStore(dst Mem, regs [fqLen]Register) {
	for i, ri := range regs {
		MOVQ(ri, dst.Offset(8*i))
	}
}

func fqMod(regs [fqLen]Register) [fqLen]Register {
	out := [fqLen]Register{GP64(), GP64(), GP64(), GP64(), GP64(), GP64()}
	for i, ri := range regs {
		MOVQ(ri, out[i])
	}
	q := Mem{Symbol: Symbol{Name: "q64"}, Base: StaticBase}
	SUBQ(q.Offset(0), out[0])
	for i := range q64[1:] {
		SBBQ(q.Offset((i+1)*8), out[i+1])
	}
	for i, ri := range regs {
		CMOVQCC(out[i], ri)
	}
	return out
}

func neg(src Mem) [fqLen]Register {
	regs := [fqLen]Register{GP64(), GP64(), GP64(), GP64(), GP64(), GP64()}
	q := Mem{Symbol: Symbol{Name: "q64"}, Base: StaticBase}
	for i, reg := range regs {
		MOVQ(q.Offset(i*8), reg)
	}

	SUBQ(src.Offset(0), regs[0])
	for i, reg := range regs[1:] {
		SBBQ(src.Offset((i+1)*8), reg)
	}

	return regs
}

func main() {
	Package("github.com/VideoCoin/go-bls12-381")

	TEXT("fqAdd", 0, "func(z *[6]uint64, x *[6]uint64, y *[6]uint64)")
	Doc("fqAdd sets z to the sum x+y.")
	x := Mem{Base: Load(Param("x"), GP64())}
	y := Mem{Base: Load(Param("y"), GP64())}
	regs := fqLoad(x)
	ADDQ(y.Offset(0), regs[0])
	for i, ri := range regs[1:] {
		ADCQ(y.Offset((i+1)*8), ri)
	}

	z := Mem{Base: Load(Param("z"), y.Base)}
	fqStore(z, fqMod(regs))

	RET()

	TEXT("fqNeg", 0, "func(z *[6]uint64, x *[6]uint64)")
	Doc("fqNeg sets z to -x.")
	x = Mem{Base: Load(Param("x"), GP64())}
	regs = neg(x)
	z = Mem{Base: Load(Param("z"), x.Base)}
	fqStore(z, regs)

	RET()

	TEXT("fqSub", 0, "func(z *[6]uint64, x *[6]uint64, y *[6]uint64)")
	Doc("fqSub sets z to the difference x-y.")
	y = Mem{Base: Load(Param("y"), GP64())}
	regs = neg(y)
	x = Mem{Base: Load(Param("x"), y.Base)}
	ADDQ(x.Offset(0), regs[0])
	for i, reg := range regs[1:] {
		ADCQ(x.Offset((i+1)*8), reg)
	}

	z = Mem{Base: Load(Param("z"), x.Base)}
	fqStore(z, fqMod(regs))

	RET()

	TEXT("fqMul", 0, "func(z *[6]uint64, x *[6]uint64, y *[6]uint64)")
	Doc("fqMul sets z to the product x*y.")

	RET()

	Generate()
}
