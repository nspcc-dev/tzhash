package gf127

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testCasesSplit = []struct {
	num *GF127x2
	h1  *GF127
	h2  *GF127
}{
	{&GF127x2{GF127{123, 31}, GF127{141, 9}}, &GF127{123, 31}, &GF127{141, 9}},
	{&GF127x2{GF127{maxUint64, 0}, GF127{0, maxUint64}}, &GF127{maxUint64, 0}, &GF127{0, maxUint64}},
}

func TestSplit(t *testing.T) {
	for _, tc := range testCasesSplit {
		a, b := Split(tc.num)
		require.Equal(t, tc.h1, a)
		require.Equal(t, tc.h2, b)
	}
}

func TestCombineTo(t *testing.T) {
	c := new(GF127x2)
	for _, tc := range testCasesSplit {
		CombineTo(tc.h1, tc.h2, c)
		require.Equal(t, tc.num, c)
	}
}

var testCasesMul10x2 = [][2]*GF127x2{
	{
		&GF127x2{GF127{123, 0}, GF127{123, 0}},
		&GF127x2{GF127{246, 0}, GF127{246, 0}},
	},
	{
		&GF127x2{GF127{maxUint64, 2}, GF127{0, 1}},
		&GF127x2{GF127{maxUint64 - 1, 5}, GF127{0, 2}},
	},
	{
		&GF127x2{GF127{0, maxUint64 >> 1}, GF127{maxUint64, 2}},
		&GF127x2{GF127{1 + 1<<63, maxUint64>>1 - 1}, GF127{maxUint64 - 1, 5}},
	},
}

func TestMul10x2(t *testing.T) {
	c := new(GF127x2)
	for _, tc := range testCasesMul10x2 {
		Mul10x2(tc[0], c)
		require.Equal(t, tc[1], c)
	}
}

var testCasesMul11x2 = [][2]*GF127x2{
	{
		&GF127x2{GF127{123, 0}, GF127{123, 0}},
		&GF127x2{GF127{141, 0}, GF127{141, 0}},
	},
	{
		&GF127x2{GF127{maxUint64, 2}, GF127{0, 1}},
		&GF127x2{GF127{1, 7}, GF127{0, 3}},
	},
	{
		&GF127x2{GF127{0, maxUint64 >> 1}, GF127{maxUint64, 2}},
		&GF127x2{GF127{1 + 1<<63, 1}, GF127{1, 7}},
	},
}

func TestMul11x2(t *testing.T) {
	c := new(GF127x2)
	for _, tc := range testCasesMul11x2 {
		Mul11x2(tc[0], c)
		require.Equal(t, tc[1], c)
	}
}
