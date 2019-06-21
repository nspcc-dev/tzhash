package gf127

import (
	"encoding/binary"
	"encoding/hex"
	"unsafe"
)

// GF127x2 represents a pair of elements of GF(2^127) stored together.
type GF127x2 [4]uint64

// Split returns 2 components of pair without additional allocations.
func Split(a *GF127x2) (*GF127, *GF127) {
	return (*GF127)(unsafe.Pointer(a)), (*GF127)(unsafe.Pointer(&(*a)[2]))
}

// CombineTo 2 elements of GF(2^127) to the respective components of pair.
func CombineTo(a *GF127, b *GF127, c *GF127x2) {
	c[0] = a[0]
	c[1] = a[1]
	c[2] = b[0]
	c[3] = b[1]
}

// Equal checks if both elements of GF(2^127) pair are equal.
func (a *GF127x2) Equal(b *GF127x2) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3]
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
	binary.BigEndian.PutUint64(buf, a[1])
	binary.BigEndian.PutUint64(buf[8:], a[0])
	binary.BigEndian.PutUint64(buf[16:], a[3])
	binary.BigEndian.PutUint64(buf[24:], a[2])
	return
}

// Mul10x2 sets (b1, b2) to (a1*x, a2*x)
func Mul10x2(a, b *GF127x2)

// Mul10x2 sets (b1, b2) to (a1*(x+1), a2*(x+1))
func Mul11x2(a, b *GF127x2)
