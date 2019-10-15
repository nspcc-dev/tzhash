package gf127

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/rand"
)

const (
	byteSize  = 16
	maxUint64 = ^uint64(0)
	msb64     = uint64(1) << 63
)

// GF127 represents element of GF(2^127)
type GF127 [2]uint64

// Random returns random element from GF(2^127).
// Is used mostly for testing.
func Random() *GF127 {
	return &GF127{rand.Uint64(), rand.Uint64() >> 1}
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
