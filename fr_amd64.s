// +build amd64,!generic

#define frLoad(from, a0, a1, a2, a3) \
        MOVQ 0+from, a0 \
        MOVQ 8+from, a1 \
        MOVQ 16+from, a2 \
        MOVQ 24+from, a3 

#define frStore(to, a0, a1, a2, a3) \
        MOVQ a0, 0+to \
        MOVQ a1, 8+to \
        MOVQ a2, 16+to \
        MOVQ a3, 24+to 

#define frMod(a0, a1, a2, a3, b0, b1, b2, b3) \
        MOVQ a0, b0 \
        MOVQ a1, b1 \
        MOVQ a2, b2 \
        MOVQ a3, b3 \
        SUBQ ·q64+0(SB), b0 \
        SBBQ ·q64+8(SB), b1 \
        SBBQ ·q64+16(SB), b2 \
        SBBQ ·q64+24(SB), b3 \
        \ // if b is negative then return a else return b
        CMOVQCC b0, a0 \
        CMOVQCC b1, a1 \
        CMOVQCC b2, a2 \
        CMOVQCC b3, a3 

TEXT ·frAdd(SB),0,$0-24
    MOVQ a+8(FP), DI
    MOVQ b+16(FP), SI
    frLoad(0(DI), R8, R9, R10, R11)
    ADDQ 0(SI), R8
    ADCQ 8(SI), R9
    ADCQ 16(SI), R10
    ADCQ 24(SI), R11
    frMod(R8, R9, R10, R11, R12, R13, R14, R15)
    MOVQ c+0(FP), DI
    frStore(0(DI), R8, R9, R10, R11)
    RET

TEXT ·frNeg(SB),0,$0-16
    MOVQ ·q64+0(SB), R8
    MOVQ ·q64+8(SB), R9
    MOVQ ·q64+16(SB), R10
    MOVQ ·q64+24(SB), R11
    MOVQ a+8(FP), DI
    SUBQ 0(DI), R8
    SBBQ 8(DI), R9
    SBBQ 16(DI), R10
    SBBQ 24(DI), R11
    MOVQ c+0(FP), DI
    frStore(0(DI), R8, R9, R10, R11)
    RET

TEXT ·frSub(SB),0,$0-24
    MOVQ ·q64+0(SB), R8
    MOVQ ·q64+8(SB), R9
    MOVQ ·q64+16(SB), R10
    MOVQ ·q64+24(SB), R11
    MOVQ b+16(FP), SI
    SUBQ 0(SI), R8
    SBBQ 8(SI), R9
    SBBQ 16(SI), R10
    SBBQ 24(SI), R11
    MOVQ a+8(FP), DI
    ADDQ 0(DI), R8
    ADCQ 8(DI), R9
    ADCQ 16(DI), R10
    ADCQ 24(DI), R11
    frMod(R8, R9, R10, R11, R12, R13, R14, R15)
    MOVQ c+0(FP), DI
    frStore(0(DI), R8, R9, R10, R11)
    RET


