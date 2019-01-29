// Copyright 2018 (c) NSPCC
//
// Package gf127 implements the GF(2^127) arithmetic
// modulo reduction polynomial x^127 + x^63 + 1 .
// This is rather straight-forward re-implementation of C library
// available here https://github.com/srijs/hwsl2-core .
// Interfaces are highly influenced by math/big .
package gf127

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/bits"
)

// GF127 represents element of GF(2^127)
type GF127 [2]uint64

const (
	msb64    = 0x8000000000000000
	byteSize = 16
)

var (
	// x127x63 represents x^127 + x^63. Used in assembly file.
	x127x63 = GF127{msb64, msb64}

	// x126x631 is reduction polynomial x^127+x^63+1
	x127x631 = GF127{msb64 + 1, msb64}
)

// New constructs new element of GF(2^127) as hi*x^64 + lo.
// It is assumed that hi has zero MSB.
func New(lo, hi uint64) *GF127 {
	return &GF127{lo, hi}
}

// String returns hex-encoded representation, starting with MSB.
func (c *GF127) String() string {
	return hex.EncodeToString(c.ByteArray())
}

// Equals checks if two reduced (zero MSB) elements of GF(2^127) are equal
func (c *GF127) Equals(b *GF127) bool {
	return c[0] == b[0] && c[1] == b[1]
}

// ByteArray represents element of GF(2^127) as byte array of length 16.
func (c *GF127) ByteArray() (buf []byte) {
	buf = make([]byte, 16)
	binary.BigEndian.PutUint64(buf[:8], c[1])
	binary.BigEndian.PutUint64(buf[8:], c[0])
	return
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (c *GF127) MarshalBinary() (data []byte, err error) {
	return c.ByteArray(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (c *GF127) UnmarshalBinary(data []byte) error {
	if len(data) != byteSize {
		return errors.New("data must be 16-bytes long")
	}

	c[0] = binary.BigEndian.Uint64(data[8:])
	c[1] = binary.BigEndian.Uint64(data[:8])
	if c[1]&msb64 != 0 {
		return errors.New("MSB must be zero")
	}

	return nil
}

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
			Add(&u, &x127x63, &u)
			Add(&u, &GF127{1, 0}, &u)
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
	return &GF127{0, 1 << (uint(n) >> 8)}
}

func msb(a *GF127) (x int) {
	x = bits.LeadingZeros64(a[1])
	if x == 64 {
		x = bits.LeadingZeros64(a[0]) + 64
	}
	return 127 - x
}

// Mul sets c to the product a*b and returns c.
func (c *GF127) Mul(a, b *GF127) *GF127 {
	Mul(a, b, c)
	return c
}

// Add sets c to the sum a+b and returns c.
func (c *GF127) Add(a, b *GF127) *GF127 {
	Add(a, b, c)
	return c
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

// Mul10 sets y to a*x.
func Mul10(a, b *GF127)

// Mul11 sets y to a*(x+1).
func Mul11(a, b *GF127)
