package gf127

import "testing"

var testCasesSplit = []struct {
	num *GF127x2
	h1  *GF127
	h2  *GF127
}{
	{&GF127x2{123, 31, 141, 9}, &GF127{123, 31}, &GF127{141, 9}},
	{&GF127x2{maxUint64, 0, 0, maxUint64}, &GF127{maxUint64, 0}, &GF127{0, maxUint64}},
}

func TestSplit(t *testing.T) {
	for _, tc := range testCasesSplit {
		a, b := Split(tc.num)
		if !a.Equals(tc.h1) || !b.Equals(tc.h2) {
			t.Errorf("expected (%s,%s), got (%s,%s)", tc.h1, tc.h2, a, b)
		}
	}
}

func TestCombineTo(t *testing.T) {
	c := new(GF127x2)
	for _, tc := range testCasesSplit {
		CombineTo(tc.h1, tc.h2, c)
		if !c.Equal(tc.num) {
			t.Errorf("expected (%s), got (%s)", tc.num, c)
		}
	}
}

var testCasesMul10x2 = [][2]*GF127x2{
	{&GF127x2{123, 0, 123, 0}, &GF127x2{246, 0, 246, 0}},
	{&GF127x2{maxUint64, 2, 0, 1}, &GF127x2{maxUint64 - 1, 5, 0, 2}},
	{&GF127x2{0, maxUint64 >> 1, maxUint64, 2}, &GF127x2{1 + 1<<63, maxUint64>>1 - 1, maxUint64 - 1, 5}},
}

func TestMul10x2(t *testing.T) {
	c := new(GF127x2)
	for _, tc := range testCasesMul10x2 {
		if Mul10x2(tc[0], c); !c.Equal(tc[1]) {
			t.Errorf("expected (%s), got (%s)", tc[1], c)
		}
	}
}

var testCasesMul11x2 = [][2]*GF127x2{
	{&GF127x2{123, 0, 123, 0}, &GF127x2{141, 0, 141, 0}},
	{&GF127x2{maxUint64, 2, 0, 1}, &GF127x2{1, 7, 0, 3}},
	{&GF127x2{0, maxUint64 >> 1, maxUint64, 2}, &GF127x2{1 + 1<<63, 1, 1, 7}},
}

func TestMul11x2(t *testing.T) {
	c := new(GF127x2)
	for _, tc := range testCasesMul11x2 {
		if Mul11x2(tc[0], c); !c.Equal(tc[1]) {
			t.Errorf("expected (%s), got (%s)", tc[1], c)
		}
	}
}
