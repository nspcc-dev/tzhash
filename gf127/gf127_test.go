package gf127

import (
	"testing"
)

const maxUint64 = ^uint64(0)

func TestAdd(t *testing.T) {
	var (
		a = Random()
		b = Random()
		e = &GF127{a[0] ^ b[0], a[1] ^ b[1]}
		c = new(GF127)
	)
	c.Add(a, b)
	if e[0] != c[0] || e[1] != c[1] {
		t.Errorf("expected (%s), got (%s)", e.String(), c.String())
	}
}

var testCasesMul = [][3]*GF127{
	// (x+1)*(x^63+x^62+...+1) == x^64+1
	{&GF127{3, 0}, &GF127{maxUint64, 0}, &GF127{1, 1}},

	// x^126 * x^2 == x^128 == x^64 + x
	{&GF127{0, 1 << 62}, &GF127{4, 0}, &GF127{2, 1}},

	// (x^64+x^63+1) * (x^64+x) == x^128+x^65+x^127+x^64+x^64+x == x^65+x^64+x^63+1
	{&GF127{1 + 1<<63, 1}, &GF127{2, 1}, &GF127{0x8000000000000001, 3}},
}

func TestMul(t *testing.T) {
	c := new(GF127)
	for _, tc := range testCasesMul {
		if Mul(tc[0], tc[1], c); !c.Equals(tc[2]) {
			t.Errorf("expected (%s), got (%s)", c.String(), tc[2].String())
		}
	}
}

var testCasesMul10 = [][2]*GF127{
	{&GF127{123, 0}, &GF127{246, 0}},
	{&GF127{maxUint64, 2}, &GF127{maxUint64 - 1, 5}},
	{&GF127{0, maxUint64 >> 1}, &GF127{1 + 1<<63, maxUint64>>1 - 1}},
}

func TestMul10(t *testing.T) {
	c := new(GF127)
	for _, tc := range testCasesMul10 {
		if Mul10(tc[0], c); !c.Equals(tc[1]) {
			t.Errorf("expected (%s), got (%s)", tc[1].String(), c.String())
		}
	}
}

var testCasesMul11 = [][2]*GF127{
	{&GF127{123, 0}, &GF127{141, 0}},
	{&GF127{maxUint64, 2}, &GF127{1, 7}},
	{&GF127{0, maxUint64 >> 1}, &GF127{1 + 1<<63, 1}},
}

func TestMul11(t *testing.T) {
	c := new(GF127)
	for _, tc := range testCasesMul11 {
		if Mul11(tc[0], c); !c.Equals(tc[1]) {
			t.Errorf("expected (%s), got (%s)", tc[1].String(), c.String())
		}
	}
}

var testCasesInv = [][2]*GF127{
	{&GF127{1, 0}, &GF127{1, 0}},
	{&GF127{54321, 12345}, &GF127{8230555108620784737, 3929873967650665114}},
}

func TestInv(t *testing.T) {
	var a, b, c = new(GF127), new(GF127), new(GF127)
	for _, tc := range testCasesInv {
		if Inv(tc[0], c); !c.Equals(tc[1]) {
			t.Errorf("expected (%s), got (%s)", tc[1].String(), c.String())
		}
	}

	for i := 0; i < 3; i++ {
		// 0 has no inverse
		if a = Random(); a.Equals(&GF127{0, 0}) {
			continue
		}
		Inv(a, b)
		Mul(a, b, c)
		if !c.Equals(&GF127{1, 0}) {
			t.Errorf("expected inverse of (%s), got (%s)", a.String(), b.String())
		}
	}
}
