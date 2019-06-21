#include "textflag.h"

// func Add(a, b, c *[2]uint64)
TEXT ·Add(SB),NOSPLIT,$0
    MOVQ a+0(FP), AX
    MOVUPD (AX), X0
    MOVQ b+8(FP), BX
    MOVUPD (BX), X1
    XORPD X1, X0
    MOVQ c+16(FP), CX
    MOVUPD X0, (CX)
    RET

// func Mul10(a, b *[2]uint64)
TEXT ·Mul10(SB),NOSPLIT,$0
    MOVQ a+0(FP), AX
    MOVUPD (AX), X0
    VPSLLQ $1, X0, X1
    VPALIGNR $8, X1, X0, X2
    PSRLQ $63, X2
    MOVUPD ·x127x63(SB), X3
    ANDPD X1, X3
    VPUNPCKHQDQ X3, X3, X3
    XORPD X2, X1
    XORPD X3, X1
    MOVQ b+8(FP), AX
    MOVUPD X1, (AX)
    RET

// func Mul10x2(a, b) *[4]uint64
TEXT ·Mul10x2(SB),NOSPLIT,$0
    MOVQ a+0(FP), AX
    VMOVDQA (AX), Y0
    VPSLLQ $1, Y0, Y1
    VPALIGNR $8, Y1, Y0, Y2
    VPSRLQ $63, Y2, Y2
    VPXOR Y1, Y2, Y2
    VPSRLQ $63, Y1, Y3
    VPSLLQ $63, Y3, Y3
    VPUNPCKHQDQ Y3, Y3, Y3
    VPXOR Y2, Y3, Y3
    MOVQ b+8(FP), AX
    VMOVDQA Y3, (AX)
    RET

// func Mul11(a, b *[2]uint64)
TEXT ·Mul11(SB),NOSPLIT,$0
    MOVQ a+0(FP), AX
    MOVUPD (AX), X0
    VPSLLQ $1, X0, X1
    VPALIGNR $8, X1, X0, X2
    PSRLQ $63, X2
    MOVUPD ·x127x63(SB), X3
    ANDPD X1, X3
    VPUNPCKHQDQ X3, X3, X3
    XORPD X2, X1
    XORPD X3, X1
    XORPD X0, X1
    MOVQ b+8(FP), AX
    MOVUPD X1, (AX)
    RET

// func Mul11x2(a, b) *[4]uint64
TEXT ·Mul11x2(SB),NOSPLIT,$0
    MOVQ a+0(FP), AX
    VMOVDQA (AX), Y0
    VPSLLQ $1, Y0, Y1
    VPALIGNR $8, Y1, Y0, Y2
    VPSRLQ $63, Y2, Y2
    VPXOR Y1, Y2, Y2
    VPSRLQ $63, Y1, Y3
    VPSLLQ $63, Y3, Y3
    VPUNPCKHQDQ Y3, Y3, Y3
    VPXOR Y2, Y3, Y3
    VPXOR Y0, Y3, Y3
    MOVQ b+8(FP), AX
    VMOVDQA Y3, (AX)
    RET

// func Mul(a, b, c *[2]uint64)
TEXT ·Mul(SB),NOSPLIT,$0
    MOVQ a+0(FP), AX              // X0 = a0 . a1
    MOVUPD (AX), X0               // X0 = a0 . a1
    MOVQ b+8(FP), BX              // X1 = b0 . b1
    MOVUPD (BX), X1               // X1 = b0 . b1
    VPUNPCKLQDQ X1, X0, X2        // X2 = a0 . b0
    VPUNPCKHQDQ X1, X0, X3        // X3 = a1 . b1
    XORPD X2, X3                  // X3 = (a0 + a1) . (b0 + b1)
    PCLMULQDQ $0x10, X3, X3       // X3 = (a0 + a1) * (b0 + b1)
    VPCLMULQDQ $0x00, X0, X1, X4  // X4 = a0 * b0
    VPCLMULQDQ $0x11, X0, X1, X5  // X5 = a1 * b1
    XORPD X4, X3                  //
    XORPD X5, X3                  // X3 = a0 * b1 + a1 * b0
    VPSLLDQ $8, X3, X2            //
    XORPD X2, X4                  // X4 = a0 * b0 + lo(X3)
    VPSRLDQ $8, X3, X6            //
    XORPD X6, X5                  // X5 = a1 * b1 + hi(X3)

    // at this point, a * b = X4 . X5 (as 256-bit number)
    // reduction modulo x^127 + x^63 + 1
    VPALIGNR $8, X4, X5, X3
    XORPD X5, X3
    PSLLQ $1, X5
    XORPD X5, X4
    VPUNPCKHQDQ X3, X5, X5
    XORPD X5, X4
    PSRLQ $63, X3
    XORPD X3, X4
    VPUNPCKLQDQ X3, X3, X5
    PSLLQ $63, X5
    XORPD X5, X4
    MOVQ c+16(FP), CX
    MOVUPD X4, (CX)
    RET
