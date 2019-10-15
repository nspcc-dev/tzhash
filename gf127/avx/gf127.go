// Copyright 2018 (c) NSPCC
//
// Package gf127 implements the GF(2^127) arithmetic
// modulo reduction polynomial x^127 + x^63 + 1 .
// This is rather straight-forward re-implementation of C library
// available here https://github.com/srijs/hwsl2-core .
// Interfaces are highly influenced by math/big .
package avx

import (
	"github.com/nspcc-dev/tzhash/gf127"
)

// GF127 is an alias for a main type.
type GF127 = gf127.GF127

const msb64 = uint64(1) << 63

var (
	// x127x63 represents x^127 + x^63. Used in assembly file.
	x127x63 = GF127{msb64, msb64}
)

// Add sets c to a+b.
func Add(a, b, c *GF127)

// Mul sets c to a*b.
func Mul(a, b, c *GF127)

// Mul10 sets b to a*x.
func Mul10(a, b *GF127)

// Mul11 sets b to a*(x+1).
func Mul11(a, b *GF127)
