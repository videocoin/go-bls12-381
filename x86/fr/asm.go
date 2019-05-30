package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

const frLen = 4

var r64 = [4]uint64{0xFFFFFFFF00000001, 0x53BDA402FFFE5BFE, 0x3339D80809A1D805, 0x73EDA753299D7D48}

func frLoad(src Mem) [frLen]Register {
	regs := [frLen]Register{GP64(), GP64(), GP64(), GP64()}
	for i, ri := range regs {
		MOVQ(src.Offset(8*i), ri)
	}
	return regs
}

func frStore(dst Mem, regs [frLen]Register) {
	for i, ri := range regs {
		MOVQ(ri, dst.Offset(8*i))
	}
}

func frMod(regs [frLen]Register) [frLen]Register {
	out := [frLen]Register{GP64(), GP64(), GP64(), GP64()}
	for i, ri := range regs {
		MOVQ(ri, out[i])
	}
	r := Mem{Symbol: Symbol{Name: "r64"}, Base: StaticBase}
	SUBQ(r.Offset(0), out[0])
	for i := range r64[1:] {
		SBBQ(r.Offset((i+1)*8), out[i+1])
	}
	for i, ri := range regs {
		CMOVQCC(out[i], ri)
	}
	return out
}

func neg(src Mem) [frLen]Register {
	regs := [frLen]Register{GP64(), GP64(), GP64(), GP64()}
	r := Mem{Symbol: Symbol{Name: "r64"}, Base: StaticBase}
	for i, reg := range regs {
		MOVQ(r.Offset(i*8), reg)
	}

	SUBQ(src.Offset(0), regs[0])
	for i, reg := range regs[1:] {
		SBBQ(src.Offset((i+1)*8), reg)
	}

	return regs
}

func main() {
	TEXT("frAdd", 0, "func(z *[4]uint64, x *[4]uint64, y *[4]uint64)")
	Doc("frAdd sets z to the sum x+y.")
	x := Mem{Base: Load(Param("x"), GP64())}
	y := Mem{Base: Load(Param("y"), GP64())}
	regs := frLoad(x)
	ADDQ(y.Offset(0), regs[0])
	for i, ri := range regs[1:] {
		ADCQ(y.Offset((i+1)*8), ri)
	}
	z := Mem{Base: Load(Param("z"), y.Base)}
	frStore(z, frMod(regs))
	RET()

	TEXT("frNeg", 0, "func(z *[4]uint64, x *[4]uint64)")
	Doc("frNeg sets z to -x.")
	x = Mem{Base: Load(Param("x"), GP64())}
	regs = neg(x)
	z = Mem{Base: Load(Param("z"), x.Base)}
	frStore(z, regs)

	TEXT("frSub", 0, "func(z *[4]uint64, x *[4]uint64, y *[4]uint64)")
	Doc("frSub sets z to the difference x-y.")
	y = Mem{Base: Load(Param("y"), GP64())}
	regs = neg(y)
	x = Mem{Base: Load(Param("x"), y.Base)}
	ADDQ(x.Offset(0), regs[0])
	for i, reg := range regs[1:] {
		ADCQ(x.Offset((i+1)*8), reg)
	}

	z = Mem{Base: Load(Param("z"), x.Base)}
	frStore(z, frMod(regs))

	RET()

	TEXT("frMul", 0, "func(z *[4]uint64, x *[4]uint64, y *[4]uint64)")
	Doc("frMul sets z to the product x*y.")

	RET()

	Generate()
}
