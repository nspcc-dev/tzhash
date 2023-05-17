package gf127

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/bits"
	"math/rand"
)

// GF127 represents element of GF(2^127).
type GF127 [2]uint64

const (
	byteSize  = 16
	maxUint64 = ^uint64(0)
	msb64     = uint64(1) << 63
)

// x127x631 is reduction polynomial x^127 + x^63 + 1.
var x127x631 = GF127{msb64 + 1, msb64}

// New constructs new element of GF(2^127) as hi*x^64 + lo.
// It is assumed that hi has zero MSB.
func New(lo, hi uint64) *GF127 {
	return &GF127{lo, hi}
}

func addGeneric(a, b, c *GF127) {
	c[0] = a[0] ^ b[0]
	c[1] = a[1] ^ b[1]
}

func mulGeneric(a, b, c *GF127) {
	r := new(GF127)
	d := *a
	for i := uint(0); i < 64; i++ {
		if b[0]&(1<<i) != 0 {
			addGeneric(r, &d, r)
		}
		mul10Generic(&d, &d)
	}
	for i := uint(0); i < 63; i++ {
		if b[1]&(1<<i) != 0 {
			addGeneric(r, &d, r)
		}
		mul10Generic(&d, &d)
	}
	*c = *r
}

func mul10Generic(a, b *GF127) {
	c := a[0] >> 63
	b[0] = a[0] << 1
	b[1] = (a[1] << 1) ^ c

	mask := b[1] & msb64
	b[0] ^= mask | (mask >> 63)
	b[1] ^= mask
}

func mul11Generic(a, b *GF127) {
	c := a[0] >> 63
	b[0] = a[0] ^ (a[0] << 1)
	b[1] = a[1] ^ (a[1] << 1) ^ c

	mask := b[1] & msb64
	b[0] ^= mask | (mask >> 63)
	b[1] ^= mask
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

		// becasuse mulAVX performs reduction on t, we need
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

// Mul1 copies b into a.
func Mul1(a, b *GF127) {
	a[0] = b[0]
	a[1] = b[1]
}

// And sets c to a & b (bitwise-and).
func And(a, b, c *GF127) {
	c[0] = a[0] & b[0]
	c[1] = a[1] & b[1]
}

// Random returns random element from GF(2^127).
// Is used mostly for testing.
func Random() *GF127 {
	return &GF127{rand.Uint64(), rand.Uint64() >> 1}
}

// String returns hex-encoded representation, starting with MSB.
func (c *GF127) String() string {
	buf := c.Bytes()
	return hex.EncodeToString(buf[:])
}

// Equals checks if two reduced (zero MSB) elements of GF(2^127) are equal.
func (c *GF127) Equals(b *GF127) bool {
	return c[0] == b[0] && c[1] == b[1]
}

// Bytes represents element of GF(2^127) as byte array of length 16.
func (c *GF127) Bytes() [16]byte {
	var buf [16]byte
	binary.BigEndian.PutUint64(buf[:8], c[1])
	binary.BigEndian.PutUint64(buf[8:], c[0])
	return buf
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (c *GF127) MarshalBinary() (data []byte, err error) {
	buf := c.Bytes()
	return buf[:], nil
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
