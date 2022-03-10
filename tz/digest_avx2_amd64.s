#include "textflag.h"

#define mulBit(bit, in_1, in_2, out_1, out_2) \
	VPSLLW      bit, Y10, Y11    \
	VPSLLQ      $1, in_1, Y1     \
	VPSRAW      $15, Y11, Y12    \
	VPALIGNR    $8, Y1, in_1, Y2 \
	VPAND       Y1, Y14, Y3      \
	VPSRLQ      $63, Y2, Y2      \
	VPUNPCKHQDQ Y3, Y3, Y3       \
	VPXOR       Y1, Y2, Y7       \
	VPXOR       Y3, in_2, out_1  \
	VPXOR       Y7, out_1, out_1 \
	VPAND       out_1, Y12, Y4   \
	VPXOR       Y4, in_1, out_2  \

// func mulByteSliceRightx2(c00c10, c01c11 *[4]uint64, n int, data *byte)
TEXT ·mulByteSliceRightx2(SB), NOSPLIT, $0
	MOVQ c00c10+0(FP), AX
	MOVQ c01c11+8(FP), BX

	VPXOR    Y13, Y13, Y13 // Y13 = 0x0000...
	VPCMPEQB Y14, Y14, Y14 // Y14 = 0xFFFF...
	VPSUBQ   Y14, Y13, Y10
	VPSLLQ   $63, Y10, Y14 // Y14 = 0x10000000... (packed quad-words with HSB set)

	MOVQ n+16(FP), CX
	MOVQ data+24(FP), DX

	VMOVDQU (AX), Y0
	VMOVDQU (BX), Y8

loop:
	CMPQ CX, $0
	JEQ  finish

	VPBROADCASTB (DX), Y10
	ADDQ         $1, DX
	SUBQ         $1, CX

	mulBit($8, Y0, Y8, Y5, Y6)
	mulBit($9, Y5, Y6, Y0, Y8)
	mulBit($10, Y0, Y8, Y5, Y6)
	mulBit($11, Y5, Y6, Y0, Y8)
	mulBit($12, Y0, Y8, Y5, Y6)
	mulBit($13, Y5, Y6, Y0, Y8)
	mulBit($14, Y0, Y8, Y5, Y6)
	mulBit($15, Y5, Y6, Y0, Y8)

	JMP loop

finish:
	VMOVDQU Y0, (AX)
	VMOVDQU Y8, (BX)

	RET
