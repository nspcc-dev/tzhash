#include "textflag.h"

// func Mul10x2(a, b) *[4]uint64
TEXT ·mul10x2AVX2(SB), NOSPLIT, $0
	MOVQ        a+0(FP), AX
	VMOVDQA     (AX), Y0
	VPSLLQ      $1, Y0, Y1
	VPALIGNR    $8, Y1, Y0, Y2
	VPSRLQ      $63, Y2, Y2
	VPXOR       Y1, Y2, Y2
	VPSRLQ      $63, Y1, Y3
	VPSLLQ      $63, Y3, Y3
	VPUNPCKHQDQ Y3, Y3, Y3
	VPXOR       Y2, Y3, Y3
	MOVQ        b+8(FP), AX
	VMOVDQA     Y3, (AX)
	RET

// func Mul11x2(a, b) *[4]uint64
TEXT ·mul11x2AVX2(SB), NOSPLIT, $0
	MOVQ        a+0(FP), AX
	VMOVDQA     (AX), Y0
	VPSLLQ      $1, Y0, Y1
	VPALIGNR    $8, Y1, Y0, Y2
	VPSRLQ      $63, Y2, Y2
	VPXOR       Y1, Y2, Y2
	VPSRLQ      $63, Y1, Y3
	VPSLLQ      $63, Y3, Y3
	VPUNPCKHQDQ Y3, Y3, Y3
	VPXOR       Y2, Y3, Y3
	VPXOR       Y0, Y3, Y3
	MOVQ        b+8(FP), AX
	VMOVDQA     Y3, (AX)
	RET
