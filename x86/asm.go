package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

// add notes about carry being always 0

const (
	fqLen        = 6
	qK64  uint64 = 0x89f3fffcfffcfffd
)

var (
	q64  = [6]uint64{0xB9FEFFFFFFFFAAAB, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}
	zero = Imm(0)
)

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
	fq := [fqLen]Register{GP64(), GP64(), GP64(), GP64(), GP64(), GP64()}
	for i, ri := range regs {
		MOVQ(ri, fq[i])
	}

	q := Mem{Symbol: Symbol{Name: "·q64"}, Base: StaticBase}
	SUBQ(q.Offset(0), fq[0])
	for i := range q64[1:] {
		SBBQ(q.Offset((i+1)*8), fq[i+1])
	}

	for i, ri := range regs {
		CMOVQCC(fq[i], ri)
	}

	return regs
}

func fqNeg(src Mem) [fqLen]Register {
	regs := [fqLen]Register{GP64(), GP64(), GP64(), GP64(), GP64(), GP64()}
	q := Mem{Symbol: Symbol{Name: "·q64"}, Base: StaticBase}
	for i, reg := range regs {
		MOVQ(q.Offset(i*8), reg)
	}

	SUBQ(src.Offset(0), regs[0])
	for i, reg := range regs[1:] {
		SBBQ(src.Offset((i+1)*8), reg)
	}

	return fqMod(regs)
}

func basicMul(z, x, y Mem) {
	product := [fqLen]Register{GP64(), GP64(), GP64(), GP64(), GP64(), GP64()}
	carry := GP64()
	for i := 0; i < fqLen; i++ {
		xi := x.Offset(i * 8)
		j := 0
		for ; j < (fqLen - 1); j++ {
			MOVQ(xi, RAX)
			MULQ(y.Offset(j * 8))
			if j == 0 {
				MOVQ(RAX, product[0])
			} else {
				ADDQ(RAX, product[j])
				ADCQ(zero, RDX)
			}
			MOVQ(RDX, product[j+1])
		}
		MOVQ(xi, RAX)
		MULQ(y.Offset(j * 8))
		ADDQ(RAX, product[j])
		ADCQ(zero, RDX)

		if i == 0 {
			MOVQ(RDX, z.Offset(fqLen*8))
		} else {
			ADDQ(z.Offset(i*8), product[0])
			for j, word := range product[1:] {
				ADCQ(z.Offset(i*8+(j+1)*8), word)
			}
			ADCQ(zero, carry)
			MOVQ(carry, z.Offset(fqLen*8+i*8))
		}
		fqStore(z.Offset(i*8), product)
	}
}

/*
func fqREDC(x Mem) [fqLen]Register {
	q := Mem{Symbol: Symbol{Name: "q64"}, Base: StaticBase}
	product := [fqLen]Register{GP64(), GP64(), GP64(), GP64(), GP64(), GP64()}
	carrySum, carryMul := GP64(), GP64()
	xi := GP64()
	for i := 0; i < fqLen; i++ {
		MOVQ(x.Offset(i*8), xi)
		MULQ(Imm(qK64))
		j := 0
		for ; j < (fqLen - 1); j++ {
			MOVQ(xi, RAX)
			MULQ(q.Offset(j * 8))
			if j == 0 {
				MOVQ(RAX, product[0])
			} else {
				ADDQ(RAX, product[j])
				ADCQ(zero, RDX)
			}
			MOVQ(RDX, product[j+1])
		}

		MOVQ(xi, RAX)
		MULQ(q.Offset(j * 8))
		ADDQ(RAX, product[j])
		ADCQ(zero, RDX)

		if i == 0 {
			MOVQ(RDX, x.Offset(fqLen*8))
		} else {
			ADDQ(x.Offset(i*8), product[0])
			for j, word := range product[1:] {
				ADCQ(x.Offset(i*8+(j+1)*8), word)
			}
			ADCQ(zero, carryMul)
			MOVQ(carryMul, x.Offset(fqLen*8+i*8))
		}
		fqStore(x.Offset(i*8), product)
	}

	return product
}
*/

func main() {
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
	x = Mem{Base: Load(Param("x"), RDI)}
	negX := fqNeg(x)
	z = Mem{Base: Load(Param("z"), x.Base)}
	fqStore(z, negX)

	RET()

	TEXT("fqSub", 0, "func(z *[6]uint64, x *[6]uint64, y *[6]uint64)")
	Doc("fqSub sets z to the difference x-y.")
	y = Mem{Base: Load(Param("y"), GP64())}
	regs = fqNeg(y)
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
	x = Mem{Base: Load(Param("x"), GP64())}
	y = Mem{Base: Load(Param("y"), GP64())}
	fqLarge := AllocLocal(96)
	basicMul(fqLarge, x, y)
	z = Mem{Base: Load(Param("z"), x.Base)}
	//fqStore(z, fqReduce(fqLarge))

	RET()

	ConstraintExpr("amd64,!generic")

	Generate()
}
