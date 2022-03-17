package gf127

import (
	"encoding/binary"
	"encoding/hex"
)

// GF127x2 represents a pair of elements of GF(2^127) stored together.
type GF127x2 [2]GF127

func mul10x2Generic(a, b *GF127x2) {
	mul10Generic(&a[0], &b[0])
	mul10Generic(&a[1], &b[1])
}

func mul11x2Generic(a, b *GF127x2) {
	mul11Generic(&a[0], &b[0])
	mul11Generic(&a[1], &b[1])
}

// Split returns 2 components of pair without additional allocations.
func Split(a *GF127x2) (*GF127, *GF127) {
	return &a[0], &a[1]
}

// CombineTo 2 elements of GF(2^127) to the respective components of pair.
func CombineTo(a *GF127, b *GF127, c *GF127x2) {
	c[0] = *a
	c[1] = *b
}

// Equal checks if both elements of GF(2^127) pair are equal.
func (a *GF127x2) Equal(b *GF127x2) bool {
	return a[0] == b[0] && a[1] == b[1]
}

// String returns hex-encoded representation, starting with MSB.
// Elements of pair are separated by comma.
func (a *GF127x2) String() string {
	b := a.ByteArray()
	return hex.EncodeToString(b[:16]) + " , " + hex.EncodeToString(b[16:])
}

// ByteArray represents element of GF(2^127) as byte array of length 32.
func (a *GF127x2) ByteArray() (buf []byte) {
	buf = make([]byte, 32)
	binary.BigEndian.PutUint64(buf, a[0][1])
	binary.BigEndian.PutUint64(buf[8:], a[0][0])
	binary.BigEndian.PutUint64(buf[16:], a[1][1])
	binary.BigEndian.PutUint64(buf[24:], a[1][0])
	return
}
