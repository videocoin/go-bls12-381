// +build arm64,!generic

// func fqAdd(z *[6]uint64, x *[6]uint64, y *[6]uint64)
TEXT ·fqAdd(SB), $0-24
    MOVD x+8(FP), R0
    MOVD (R0), R1
    MOVD 8(R0), R2
    MOVD 16(R0), R3
    MOVD 24(R0), R4
    MOVD 32(R0), R5
    MOVD 40(R0), R6

    MOVD(y+16), R0
    MOVD (R0), R7
    MOVD 8(R0), R8
    MOVD 16(R0), R9
    MOVD 24(R0), R10
    MOVD 32(R0), R11
    MOVD 40(R0), R12

    ADDS R7, R1
    ADCS R8, R2
    ADCS R9, R3
    ADCS R10, R4
    ADCS R11, R5
    ADCS R12, R6

    MOVD ZR, R0
    SUBS ·q64+0(SB), R1, R7
	SBCQ ·q64+8(SB), R2, R8
	SBCQ ·q64+16(SB), R3, R9
	SBCQ ·q64+24(SB), R4, R10
	SBCQ ·q64+32(SB), R5, R11
	SBCQ ·q64+40(SB), R6, R12
    SBCS  ZR, R0, R0

    CSEL CS, R7, R1, R1
    CSEL CS, R8, R2, R2
    CSEL CS, R9, R3, R3
    CSEL CS, R10, R4, R4
    CSEL CS, R11, R5, R5
    CSEL CS, R12, R6, R6

    MOVD z+0(FP), R0
    MOVD R1, (R0)
    MOVD R2, 8(R0)
    MOVD R3, 16(R0)
    MOVD R4, 24(R0)
    MOVD R5, 32(R0)
    MOVD R6, 40(R0)
    
    RET

// func fqNeg(z *[6]uint64, x *[6]uint64)
TEXT ·fqNeg(SB), $0-16
    MOVD x+8(FP), R0
    MOVD (R0), R1
    MOVD 8(R0), R2
    MOVD 16(R0), R3
    MOVD 24(R0), R4
    MOVD 32(R0), R5
    MOVD 40(R0), R6

    MOVD ·q64+0(SB), R7
    MOVD ·q64+8(SB), R8
    MOVD ·q64+16(SB), R9
    MOVD ·q64+24(SB), R10
    MOVD ·q64+32(SB), R11
    MOVD ·q64+40(SB), R12

    SUBS R7, R1, R7
    SBCS R8, R2, R8
    SBCS R9, R3, R9
    SBCS R10, R4, R10
    SBCS R11, R5, R11
    SBCS R12, R6, R12

    CSEL CS, R7, R1, R1
	CSEL CS, R8, R2, R2
	CSEL CS, R9, R3, R3
	CSEL CS, R10, R4, R4
    CSEL CS, R11, R5, R5
    CSEL CS, R12, R6, R6

    MOVD z+0(FP), R0
    MOVD R1, (R0)
    MOVD R2, 8(R0)
    MOVD R3, 16(R0)
    MOVD R4, 24(R0)
    MOVD R5, 32(R0)
    MOVD R6, 40(R0)

    RET

// func fqSub(z *[6]uint64, x *[6]uint64, y *[6]uint64)
TEXT ·fqSub(SB), $0-24
    MOVD x+8(FP), R0
    MOVD (R0), R1
    MOVD 8(R0), R2
    MOVD 16(R0), R3
    MOVD 24(R0), R4
    MOVD 32(R0), R5
    MOVD 40(R0), R6

    MOVD(y+16), R0
    MOVD (R0), R7
    MOVD 8(R0), R8
    MOVD 16(R0), R9
    MOVD 24(R0), R10
    MOVD 32(R0), R11
    MOVD 40(R0), R12

    RET

TEXT ·fqMul(SB), $96-24
    RET