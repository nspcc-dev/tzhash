#include "textflag.h"

#define mask(bit, src, tmp, to1, to2) \
    MOVQ src, tmp \
    SHRQ bit, tmp \
    ANDQ $1, tmp  \
    NEGQ tmp      \
    MOVQ tmp, to1 \
    VPBROADCASTB to1, to2

#define mulBit(bit) \
    VPSLLQ $1, Y0, Y1 \
    VPALIGNR $8, Y1, Y0, Y2 \
    VPSRLQ $63, Y2, Y2 \
    VPXOR Y1, Y2, Y2 \
    VPSRLQ $63, Y1, Y3 \
    VPSLLQ $63, Y3, Y3 \
    VPUNPCKHQDQ Y3, Y3, Y3 \
    VPXOR Y2, Y3, Y3 \
    mask(bit, CX, DX, X1, Y2) \
    VPXOR Y3, Y8, Y3 \
    VPAND Y3, Y2, Y4 \
    VPXOR Y4, Y0, Y8 \
    VMOVDQA Y3, Y0

// func mulByteRightx2(c00c10, c01c11 *[4]uint64, b byte)
TEXT ·mulByteRightx2(SB),NOSPLIT,$0
    MOVQ c00c10+0(FP), AX
    VMOVDQA (AX), Y0
    MOVQ c01c11+8(FP), BX
    VMOVDQA (BX), Y8
    MOVB b+16(FP), CX

    mulBit($7)
    mulBit($6)
    mulBit($5)
    mulBit($4)
    mulBit($3)
    mulBit($2)
    mulBit($1)
    mulBit($0)

    VMOVDQA Y8, (BX)
    VMOVDQA Y0, (AX)

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
