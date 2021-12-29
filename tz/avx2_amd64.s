#include "textflag.h"

#define mask(bit, tmp, to) \
    VPSRLW bit, Y10, tmp \
    VPAND Y12, tmp, to \ // to2 = 0x000<bit>000<bit>...
    VPSUBW to, Y13, to  // to2 = 0xFFFF.. or 0x0000 depending on bit

#define mulBit(bit) \
    VPSLLQ $1, Y0, Y1 \
    VPALIGNR $8, Y1, Y0, Y2 \
    VPSRLQ $63, Y2, Y2 \
    VPXOR Y1, Y2, Y2 \
    VPAND Y1, Y14, Y3 \
    VPUNPCKHQDQ Y3, Y3, Y3 \
    VPXOR Y2, Y3, Y3 \
    mask(bit, Y11, Y2) \
    VPXOR Y3, Y8, Y3 \
    VPAND Y3, Y2, Y4 \
    VPXOR Y4, Y0, Y8 \
    VMOVDQU Y3, Y0

// func mulByteRightx2(c00c10, c01c11 *[4]uint64, b byte)
TEXT ·mulByteRightx2(SB),NOSPLIT,$0
    MOVQ c00c10+0(FP), AX
    VMOVDQU (AX), Y0
    MOVQ c01c11+8(FP), BX
    VMOVDQU (BX), Y8

    VPXOR Y13, Y13, Y13    // Y13 = 0x0000...
    VPCMPEQB Y14, Y14, Y14 // Y14 = 0xFFFF...
    VPSUBQ Y14, Y13, Y10
    VPSUBW Y14, Y13, Y12   // Y12 = 0x00010001... (packed words of 1)
    VPSLLQ $63, Y10, Y14   // Y14 = 0x10000000... (packed quad-words with HSB set)

    VPBROADCASTB b+16(FP), X10 // X10 = packed bytes of b.
    VPMOVZXBW X10, Y10         // Extend with zeroes to packed words.

    mulBit($7)
    mulBit($6)
    mulBit($5)
    mulBit($4)
    mulBit($3)
    mulBit($2)
    mulBit($1)
    mulBit($0)

    VMOVDQU Y8, (BX)
    VMOVDQU Y0, (AX)

    RET

// func mulBitRightx2(c00c10, c01c11 *[4]uint64, e *[2]uint64)
TEXT ·mulBitRightx2(SB),NOSPLIT,$0
    MOVQ c00c10+0(FP), AX
    VMOVDQU (AX), Y0
    MOVQ c01c11+8(FP), BX
    VMOVDQU (BX), Y8

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
    VMOVDQU Y8, (BX)
    VMOVDQU Y3, (AX)
    RET
