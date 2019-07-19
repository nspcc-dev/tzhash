#include "textflag.h"

// func mulBitRight(c00, c01, c10, c11, e *[2]uint64)
TEXT ·mulBitRight(SB),NOSPLIT,$0
    MOVQ c00+0(FP), AX
    MOVUPD (AX), X0
    MOVUPD X0, X8             // remember c00 value
    MOVQ c01+8(FP), BX
    MOVUPD (BX), X1
    MOVQ c10+16(FP), CX
    MOVUPD (CX), X2
    MOVUPD X2, X9             // remember c10 value
    MOVQ c11+24(FP), DX
    MOVUPD (DX), X3

    // c00 *= 2
    VPSLLQ $1, X0, X5
    VPALIGNR $8, X5, X0, X6
    PSRLQ $63, X6
    MOVUPD ·x127x63(SB), X7
    ANDPD X5, X7
    VPUNPCKHQDQ X7, X7, X7
    XORPD X6, X5
    XORPD X7, X5
    MOVUPD X5, X0

    // c00 += c01
    XORPD X1, X0
    MOVUPD X0, (AX)

    // c10 *= 2
    VPSLLQ $1, X2, X5
    VPALIGNR $8, X5, X2, X6
    PSRLQ $63, X6
    MOVUPD ·x127x63(SB), X7
    ANDPD X5, X7
    VPUNPCKHQDQ X7, X7, X7
    XORPD X6, X5
    XORPD X7, X5
    MOVUPD X5, X2

    // c10 += c11
    XORPD X3, X2
    MOVUPD X2, (CX)

    MOVQ e+32(FP), AX
    MOVUPD (AX), X5

    // c01 = c00 + e
    VANDPD X0, X5, X1

    // c01 += X8 (old c00)
    XORPD X8, X1
    MOVUPD X1, (BX)

    // c11 = c10 + e
    VANDPD X2, X5, X3

    // c11 += X9 (old c10)
    XORPD X9, X3
    MOVUPD X3, (DX)
    RET


// func mulBitRightx2(c00c10, c01c11 *[4]uint64, e *[2]uint64)
TEXT ·mulBitRightx2(SB),NOSPLIT,$0
    MOVQ c00c10+0(FP), AX
    VMOVDQA (AX), Y0
    MOVQ c01c11+8(FP), BX
    VMOVDQA (BX), Y8

    VPSLLQ $1, Y0, Y1
    VPALIGNR $8, Y1, Y0, Y2
    VPSRLQ $63, Y2, Y2
    VPXOR Y1, Y2, Y2
    VPSRLQ $63, Y1, Y3
    VPSLLQ $63, Y3, Y3
    VPUNPCKHQDQ Y3, Y3, Y3
    VPXOR Y2, Y3, Y3

    MOVQ e+16(FP), CX
    VBROADCASTI128 (CX), Y2

    VPXOR Y3, Y8, Y3
    VPAND Y3, Y2, Y4
    VPXOR Y4, Y0, Y8
    VMOVDQA Y8, (BX)
    VMOVDQA Y3, (AX)
    RET
