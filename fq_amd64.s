// +build amd64,!generic

#define fqLoad(from, a0, a1, a2, a3, a4, a5) \
        MOVQ 0+from, a0 \
        MOVQ 8+from, a1 \
        MOVQ 16+from, a2 \
        MOVQ 24+from, a3 \
        MOVQ 32+from, a4 \
        MOVQ 40+from, a5 

#define fqStore(to, a0, a1, a2, a3, a4, a5) \
        MOVQ a0, 0+to \
        MOVQ a1, 8+to \
        MOVQ a2, 16+to \
        MOVQ a3, 24+to \
        MOVQ a4, 32+to \
        MOVQ a5, 40+to

#define fqMod(a0, a1, a2, a3, a4, a5, b0, b1, b2, b3, b4, b5) \
        MOVQ a0, b0 \
        MOVQ a1, b1 \
        MOVQ a2, b2 \
        MOVQ a3, b3 \
        MOVQ a4, b4 \
        MOVQ a5, b5 \
        SUBQ ·q64+0(SB), b0 \
        SBBQ ·q64+8(SB), b1 \
        SBBQ ·q64+16(SB), b2 \
        SBBQ ·q64+24(SB), b3 \
        SBBQ ·q64+32(SB), b4 \
        SBBQ ·q64+40(SB), b5 \
        \ // if b is negative then return a else return b
        CMOVQCC b0, a0 \
        CMOVQCC b1, a1 \
        CMOVQCC b2, a2 \
        CMOVQCC b3, a3 \
        CMOVQCC b4, a4 \
        CMOVQCC b5, a5
          
TEXT ·fqAdd(SB),0,$0-24
    MOVQ a+8(FP), DI
    MOVQ b+16(FP), SI
    fqLoad(0(DI), R8, R9, R10, R11, R12, R13)
    ADDQ 0(SI), R8
    ADCQ 8(SI), R9
    ADCQ 16(SI), R10
    ADCQ 24(SI), R11
    ADCQ 32(SI), R12
    ADCQ 40(SI), R13
    fqMod(R8, R9, R10, R11, R12, R13, R14, R15, AX, BX, CX, DX)
    MOVQ c+0(FP), DI
    fqStore(0(DI), R8, R9, R10, R11, R12, R13)
    RET

TEXT ·fqNeg(SB),0,$0-16
    MOVQ ·q64+0(SB), R8
    MOVQ ·q64+8(SB), R9
    MOVQ ·q64+16(SB), R10
    MOVQ ·q64+24(SB), R11
    MOVQ ·q64+32(SB), R12
    MOVQ ·q64+40(SB), R13
    MOVQ a+8(FP), DI
    SUBQ 0(DI), R8
    SBBQ 8(DI), R9
    SBBQ 16(DI), R10
    SBBQ 24(DI), R11
    SBBQ 32(DI), R12
    SBBQ 40(DI), R13
    MOVQ c+0(FP), DI
    fqStore(0(DI), R8, R9, R10, R11, R12, R13)
    RET

TEXT ·fqSub(SB),0,$0-24
    MOVQ ·q64+0(SB), R8
    MOVQ ·q64+8(SB), R9
    MOVQ ·q64+16(SB), R10
    MOVQ ·q64+24(SB), R11
    MOVQ ·q64+32(SB), R12
    MOVQ ·q64+40(SB), R13
    MOVQ b+16(FP), SI
    SUBQ 0(SI), R8
    SBBQ 8(SI), R9
    SBBQ 16(SI), R10
    SBBQ 24(SI), R11
    SBBQ 32(SI), R12
    SBBQ 40(SI), R13
    MOVQ a+8(FP), DI
    ADDQ 0(DI), R8
    ADCQ 8(DI), R9
    ADCQ 16(DI), R10
    ADCQ 24(DI), R11
    ADCQ 32(DI), R12
    ADCQ 40(DI), R13
    fqMod(R8, R9, R10, R11, R12, R13, R14, R15, AX, BX, CX, DX)
    MOVQ c+0(FP), DI
    fqStore(0(DI), R8, R9, R10, R11, R12, R13)
    RET

TEXT ·fqBasicMul(SB),0,$0-24
    RET

TEXT ·fqMul(SB),0,$0-24
    RET
    