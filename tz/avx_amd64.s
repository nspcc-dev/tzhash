#include "textflag.h"

// mul2 multiplicates FROM by 2, stores result in R1
// and uses R1, R2 and R3 for internal computations.
#define mul2(FROM, TO, R2, R3) \
    VPSLLQ $1, FROM, TO \
    VPALIGNR $8, TO, FROM, R2 \
    VPSRLQ $63, R2, R2 \
    VMOVDQU ·x127x63(SB), R3 \
    VANDPD TO, R3, R3 \
    VPUNPCKHQDQ R3, R3, R3 \
    VXORPD R2, TO, TO \
    VXORPD R3, TO, TO

#define mask(bit, src, tmp, to1, to2) \
    MOVQ src, tmp \
    SHRQ bit, tmp \
    ANDQ $1, tmp  \
    NEGQ tmp      \
    MOVQ tmp, to1 \
    VSHUFPS $0, to1, to1, to2
    // VPBROADCASTB to1, to2
    // Can't use VPBROADCASTB because it is AVX2 instruction
    //https://software.intel.com/en-us/forums/intel-isa-extensions/topic/301461

#define mulBit(bit) \
    VMOVDQU X0, X8 \
    VMOVDQU X2, X9 \
    mul2(X0, X5, X6, X7) \
    VXORPD X1, X5, X0 \
    mul2(X2, X5, X6, X7) \
    VXORPD X3, X5, X2 \
    mask(bit, CX, DX, X6, X5) \
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
    MOVB b+32(FP), CX

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
