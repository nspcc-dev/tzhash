#include "textflag.h"

// mul2 multiplicates FROM by 2, stores result in R1
// and uses R1, R2 and R3 for internal computations.
#define mul2(FROM, TO, R2, R3) \
    VPSLLQ $1, FROM, TO \
    VPALIGNR $8, TO, FROM, R2 \
    VPSRLQ $63, R2, R2 \
    VANDPD TO, X14, R3 \
    VPUNPCKHQDQ R3, R3, R3 \
    VXORPD R2, TO, TO \
    VXORPD R3, TO, TO

#define mask(bit, tmp, to) \
   VPSRLW bit, X10, tmp \
   VPAND X12, tmp, to \ // to = 0x000<bit>000<bit>...
   VPSUBW to, X13, to   // to = 0xFFFF.. or 0x0000 depending on bit

#define mulBit(bit) \
    VMOVDQU X0, X8 \
    VMOVDQU X2, X9 \
    mul2(X0, X5, X6, X7) \
    VXORPD X1, X5, X0 \
    mul2(X2, X5, X6, X7) \
    VXORPD X3, X5, X2 \
    mask(bit, X6, X5) \
    VANDPD X0, X5, X1 \
    VXORPD X8, X1, X1 \
    VANDPD X2, X5, X3 \
    VXORPD X9, X3, X3

// func mulBitRight(c00, c01, c10, c11, e *[2]uint64)
TEXT ·mulBitRight(SB),NOSPLIT,$0
    MOVQ c00+0(FP), AX
    VMOVDQU (AX), X0
    VMOVDQU X0, X8             // remember c00 value
    MOVQ c01+8(FP), BX
    VMOVDQU (BX), X1
    MOVQ c10+16(FP), CX
    VMOVDQU (CX), X2
    VMOVDQU X2, X9             // remember c10 value
    MOVQ c11+24(FP), DX
    VMOVDQU (DX), X3

    VPXOR X13, X13, X13     // Y13 = 0x0000...
    VPCMPEQB X14, X14, X14  // Y14 = 0xFFFF...
    VPSUBQ X14, X13, X13
    VPSLLQ $63, X13, X14

    mul2(X0, X5, X6, X7) // c00 *= 2
    VXORPD X5, X1, X0    // c00 += c01
    mul2(X2, X5, X6, X7) // c10 *= 2
    VXORPD X3, X5, X2    // c10 += c11
    MOVQ e+32(FP), CX
    VMOVDQU (CX), X5
    VANDPD X0, X5, X1    // c01 = c00 + e
    VXORPD X8, X1, X1    // c01 += X8 (old c00)
    VANDPD X2, X5, X3    // c11 = c10 + e
    VXORPD X9, X3, X3         // c11 += x9 (old c10)

    VMOVDQU X0, (AX)
    MOVQ c10+16(FP), CX
    VMOVDQU X2, (CX)
    VMOVDQU X1, (BX)
    VMOVDQU X3, (DX)

    RET

TEXT ·mulByteRight(SB),NOSPLIT,$0
    MOVQ c00+0(FP), AX
    VMOVDQU (AX), X0
    MOVQ c01+8(FP), BX
    VMOVDQU (BX), X1
    MOVQ c10+16(FP), CX
    VMOVDQU (CX), X2
    MOVQ c11+24(FP), DX
    VMOVDQU (DX), X3
    MOVQ $0, CX
    MOVB b+32(FP), CX

    VPXOR X13, X13, X13     // X13 = 0x0000...
    VPCMPEQB X14, X14, X14  // X14 = 0xFFFF...
    VPSUBQ X14, X13, X10
    VPSUBW X14, X13, X12    // X12 = 0x00010001... (packed words of 1)
    VPSLLQ $63, X10, X14    // X14 = 0x10000000... (packed quad-words with HSB set)

    MOVQ CX, X10
    VPSHUFLW $0, X10, X11
    VPSHUFD $0, X11, X10

    mulBit($7)
    mulBit($6)
    mulBit($5)
    mulBit($4)
    mulBit($3)
    mulBit($2)
    mulBit($1)
    mulBit($0)

    VMOVDQU X0, (AX)
    MOVQ c10+16(FP), CX
    VMOVDQU X2, (CX)
    VMOVDQU X1, (BX)
    MOVQ c11+24(FP), DX
    VMOVDQU X3, (DX)

    RET
