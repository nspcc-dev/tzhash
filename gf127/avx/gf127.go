// Copyright 2018 (c) NSPCC
//
// Package gf127 implements the GF(2^127) arithmetic
// modulo reduction polynomial x^127 + x^63 + 1 .
// This is rather straight-forward re-implementation of C library
// available here https://github.com/srijs/hwsl2-core .
// Interfaces are highly influenced by math/big .
package avx

import (
	"math/bits"

	"github.com/nspcc-dev/tzhash/gf127"
)

// GF127 is an alias for a main type.
type GF127 = gf127.GF127

const msb64 = uint64(1) << 63

var (
	// x127x63 represents x^127 + x^63. Used in assembly file.
	x127x63 = GF127{msb64, msb64}

	// x126x631 is reduction polynomial x^127+x^63+1
	x127x631 = GF127{msb64 + 1, msb64}
)

// Inv sets b to a^-1
// Algorithm is based on Extended Euclidean Algorithm
// and is described by Hankerson, Hernandez, Menezes in
// https://link.springer.com/content/pdf/10.1007/3-540-44499-8_1.pdf
func Inv(a, b *GF127) {
	var (
		v    = x127x631
		u    = *a
		c, d = &GF127{1, 0}, &GF127{0, 0}
		t    = new(GF127)
		x    *GF127
	)

	// degree of polynomial is a position of most significant bit
	for du, dv := msb(&u), msb(&v); du != 0; du, dv = msb(&u), msb(&v) {
		if du < dv {
			v, u = u, v
			dv, du = du, dv
			d, c = c, d
		}

		x = xN(du - dv)

		Mul(x, &v, t)
		Add(&u, t, &u)

		// becasuse mul performs reduction on t, we need
		// manually reduce u at first step
		if msb(&u) == 127 {
			Add(&u, &x127x631, &u)
		}

		Mul(x, d, t)
		Add(c, t, c)
	}
	*b = *c
}

func xN(n int) *GF127 {
	if n < 64 {
		return &GF127{1 << uint(n), 0}
	}
	return &GF127{0, 1 << uint(n-64)}
}

func msb(a *GF127) (x int) {
	x = bits.LeadingZeros64(a[1])
	if x == 64 {
		x = bits.LeadingZeros64(a[0]) + 64
	}
	return 127 - x
}

// Mul1 copies a to b.
func Mul1(a, b *GF127) {
	b[0] = a[0]
	b[1] = a[1]
}

// And sets c to a & b (bitwise-and).
func And(a, b, c *GF127) {
	c[0] = a[0] & b[0]
	c[1] = a[1] & b[1]
}

// Add sets c to a+b.
func Add(a, b, c *GF127)

// Mul sets c to a*b.
func Mul(a, b, c *GF127)

// Mul10 sets b to a*x.
func Mul10(a, b *GF127)

// Mul11 sets b to a*(x+1).
func Mul11(a, b *GF127)
