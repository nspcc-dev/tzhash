package gf127

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	var (
		a = Random()
		b = Random()
		e = &GF127{a[0] ^ b[0], a[1] ^ b[1]}
		c = new(GF127)
	)
	Add(a, b, c)
	require.Equal(t, e, c)
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
		Mul(tc[0], tc[1], c)
		require.Equal(t, tc[2], c)
	}
}

func TestMulInPlace(t *testing.T) {
	for _, tc := range testCasesMul {
		a := *tc[0]
		b := *tc[1]
		Mul(&a, &b, &b)
		require.Equal(t, a, *tc[0])
		require.Equal(t, b, *tc[2])

		b = *tc[1]
		Mul(&a, &b, &a)
		require.Equal(t, b, *tc[1])
		require.Equal(t, a, *tc[2])
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
		Mul10(tc[0], c)
		require.Equal(t, tc[1], c)
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
		Mul11(tc[0], c)
		require.Equal(t, tc[1], c)
	}
}

var testCasesInv = [][2]*GF127{
	{&GF127{1, 0}, &GF127{1, 0}},
	{&GF127{3, 0}, &GF127{msb64, ^msb64}},
	{&GF127{54321, 12345}, &GF127{8230555108620784737, 3929873967650665114}},
}

func TestInv(t *testing.T) {
	var a, b, c = new(GF127), new(GF127), new(GF127)
	for _, tc := range testCasesInv {
		Inv(tc[0], c)
		require.Equal(t, tc[1], c)
	}

	for range 3 {
		// 0 has no inverse
		if a = Random(); a.Equals(&GF127{0, 0}) {
			continue
		}
		Inv(a, b)
		Mul(a, b, c)
		require.Equal(t, &GF127{1, 0}, c)
	}
}

func TestGF127_MarshalBinary(t *testing.T) {
	a := New(0xFF, 0xEE)
	data, err := a.MarshalBinary()
	require.NoError(t, err)
	require.Equal(t, data, []byte{0, 0, 0, 0, 0, 0, 0, 0xEE, 0, 0, 0, 0, 0, 0, 0, 0xFF})

	a = Random()
	data, err = a.MarshalBinary()
	require.NoError(t, err)

	b := new(GF127)
	err = b.UnmarshalBinary(data)
	require.NoError(t, err)
	require.Equal(t, a, b)

	err = b.UnmarshalBinary([]byte{0, 1, 2, 3})
	require.Error(t, err)
}
