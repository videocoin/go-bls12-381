package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

// TODO(rgeraldes): add notes about carry being always 0

const (
	fqLen = 6
	qK64  = 0x89f3fffcfffcfffd
)

var (
	zero    = Imm(0)
	hasBMI2 = Mem{Symbol: Symbol{Name: "·hasBMI2"}, Base: StaticBase}
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
	for i := 0; i < (fqLen - 1); i++ {
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

func fqMul() {
	x := Mem{Base: Load(Param("x"), GP64())}
	y := Mem{Base: Load(Param("y"), GP64())}
	fqLarge := AllocLocal(96)
	product := [fqLen]Register{GP64(), GP64(), GP64(), GP64(), GP64(), GP64()}
	carryMul, carrySum := GP64(), GP64()
	CMPB(hasBMI2, zero)
	JE(LabelRef("fallback"))
	basicMulBMI2(product, fqLarge, x, y)
	fqREDCBMI2(product, carrySum, fqLarge)
	JMP(LabelRef("out"))
	Label("fallback")
	basicMul(product, fqLarge, x, y)
	fqREDC(product, carryMul, carrySum, fqLarge)
	Label("out")
	z := Mem{Base: Load(Param("z"), x.Base)}
	fqStore(z, product)
	RET()
}

func basicMul(product [fqLen]Register, z Mem, x Mem, y Mem) {
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
			ADCQ(zero, RDX)
			MOVQ(RDX, z.Offset(fqLen*8+i*8))
		}
		fqStore(z.Offset(i*8), product)
	}
}

func basicMulBMI2(product [fqLen]Register, z Mem, x Mem, y Mem) {
	for i := 0; i < fqLen; i++ {
		MOVQ(x.Offset(i*8), RDX)
		j := 0
		for ; j < (fqLen - 1); j++ {
			if j == 0 {
				MULXQ(y.Offset(j*8), product[j], product[j+1])
			} else {
				MULXQ(y.Offset(j*8), RAX, product[j+1])
				ADDQ(RAX, product[j])
				ADCQ(zero, product[j+1])
			}
		}
		MULXQ(y.Offset(j*8), RAX, RBX)
		ADDQ(RAX, product[j])
		ADCQ(zero, RBX)

		if i == 0 {
			MOVQ(RBX, z.Offset(fqLen*8))
		} else {
			ADDQ(z.Offset(i*8), product[0])
			for j, word := range product[1:] {
				ADCQ(z.Offset(i*8+(j+1)*8), word)
			}
			ADCQ(zero, RBX)
			MOVQ(RBX, z.Offset(fqLen*8+i*8))
		}
		fqStore(z.Offset(i*8), product)
	}
}

func fqREDC(product [fqLen]Register, carryMul Register, carrySum Register, x Mem) {
	XORQ(carrySum, carrySum)
	q := Mem{Symbol: Symbol{Name: "·q64"}, Base: StaticBase}
	for i := 0; i < fqLen; i++ {
		j := 0
		for ; j < (fqLen - 1); j++ {
			MOVQ(Imm(qK64), RAX)
			MULQ(x.Offset(i * 8))
			MULQ(q.Offset(j * 8))
			if j == 0 {
				MOVQ(RAX, product[0])
			} else {
				ADDQ(RAX, product[j])
				ADCQ(zero, RDX)
			}
			MOVQ(RDX, product[j+1])
		}
		MOVQ(Imm(qK64), RAX)
		MULQ(x.Offset(i * 8))
		MULQ(q.Offset(j * 8))
		ADDQ(RAX, product[j])
		ADCQ(zero, RDX)

		MOVQ(RDX, carryMul)
		ADDQ(x.Offset(i*8), product[0])
		for j, word := range product[1:] {
			ADCQ(x.Offset(i*8+(j+1)*8), word)
		}

		ADCQ(zero, carryMul)
		ADDQ(carrySum, carryMul)
		XORQ(carrySum, carrySum)
		ADDQ(x.Offset(fqLen*8+i*8), carryMul)
		ADCQ(zero, carrySum)

		MOVQ(carryMul, x.Offset(fqLen*8+i*8))

		// note(rgeraldes): since the low order bits are going to be discarded and
		// x[i+j=0] is not used anymore during the program, we can skip the assignment.
		for j, pi := range product[1:] {
			MOVQ(pi, x.Offset(i*8+(j+1)*8))
		}
	}

	for i := 0; i < fqLen; i++ {
		MOVQ(x.Offset(fqLen*8+i*8), product[i])
	}

	fqMod(product)
}

func fqREDCBMI2(product [fqLen]Register, carrySum Register, x Mem) {
	XORQ(carrySum, carrySum)
	q := Mem{Symbol: Symbol{Name: "·q64"}, Base: StaticBase}
	for i := 0; i < fqLen; i++ {
		MOVQ(Imm(qK64), RDX)
		MULXQ(x.Offset(i*8), RDX, RAX)
		j := 0
		for ; j < (fqLen - 1); j++ {
			if j == 0 {
				MULXQ(q.Offset(j*8), product[j], product[j+1])
			} else {
				MULXQ(q.Offset(j*8), RAX, product[j+1])
				ADDQ(RAX, product[j])
				ADCQ(zero, product[j+1])
			}
		}
		MULXQ(q.Offset(j*8), RAX, RBX)
		ADDQ(RAX, product[j])
		ADCQ(zero, RBX)

		ADDQ(x.Offset(i*8), product[0])
		for j, word := range product[1:] {
			ADCQ(x.Offset(i*8+(j+1)*8), word)
		}

		ADCQ(zero, RBX)
		ADDQ(carrySum, RBX)
		XORQ(carrySum, carrySum)
		ADDQ(x.Offset(fqLen*8+i*8), RBX)
		ADCQ(zero, carrySum)

		MOVQ(RBX, x.Offset(fqLen*8+i*8))

		// note(rgeraldes): since the low order bits are going to be discarded and
		// x[i+j=0] is not used anymore during the program, we can skip the assignment.
		for j, pi := range product[1:] {
			MOVQ(pi, x.Offset(i*8+(j+1)*8))
		}
	}

	for i := 0; i < fqLen; i++ {
		MOVQ(x.Offset(fqLen*8+i*8), product[i])
	}

	fqMod(product)
}

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
	// Replace RDI with gp64()
	x = Mem{Base: Load(Param("x"), RDI)}
	negX := fqNeg(x)
	z = Mem{Base: Load(Param("z"), x.Base)}
	fqStore(z, negX)

	RET()

	/*
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
	*/

	TEXT("fqSub", 0, "func(z *[6]uint64, x *[6]uint64, y *[6]uint64)")
	Doc("fqSub sets z to the difference x-y.")
	x = Mem{Base: Load(Param("x"), GP64())}
	regs = fqLoad(x)
	y = Mem{Base: Load(Param("y"), x.Base)}
	SUBQ(y.Offset(0), regs[0])
	for i, ri := range regs[1:] {
		SBBQ(y.Offset((i+1)*8), ri)
	}
	regsQ := [fqLen]Register{GP64(), GP64(), GP64(), GP64(), GP64(), GP64()}
	q := Mem{Symbol: Symbol{Name: "·q64"}, Base: StaticBase}
	for i, ri := range regsQ {
		MOVQ(q.Offset(i*8), ri)
	}
	r0 := y.Base
	XORQ(r0, r0)
	for _, ri := range regsQ {
		CMOVQCC(r0, ri)
	}

	ADDQ(regsQ[0], regs[0])
	for i, ri := range regs[1:] {
		ADCQ(regsQ[i], ri)
	}

	z = Mem{Base: Load(Param("z"), x.Base)}
	fqStore(z, fqMod(regs))

	RET()

	TEXT("fqMul", 0, "func(z *[6]uint64, x *[6]uint64, y *[6]uint64)")
	Doc("fqMul sets z to the product x*y.")
	fqMul()

	ConstraintExpr("amd64,!generic")
	Generate()
}
