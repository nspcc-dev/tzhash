package tz

import (
	"hash"

	"github.com/nspcc-dev/tzhash/gf127"
)

type digest3 struct {
	x [2]gf127.GF127x2
}

// type assertion
var _ hash.Hash = (*digest3)(nil)

func newAVX2Inline() *digest3 {
	d := new(digest3)
	d.Reset()
	return d
}

func (d *digest3) Write(data []byte) (n int, err error) {
	n = len(data)
	if len(data) != 0 {
		mulByteSliceRightx2(&d.x[0], &d.x[1], n, &data[0])
	}
	return
}

func (d *digest3) Sum(in []byte) []byte {
	// Make a copy of d so that caller can keep writing and summing.
	d0 := *d
	h := d0.checkSum()
	return append(in, h[:]...)
}
func (d *digest3) Reset() {
	d.x[0] = gf127.GF127x2{GF127{1, 0}, GF127{0, 0}}
	d.x[1] = gf127.GF127x2{GF127{0, 0}, GF127{1, 0}}
}
func (d *digest3) Size() int      { return Size }
func (d *digest3) BlockSize() int { return hashBlockSize }
func (d *digest3) checkSum() (b [Size]byte) {
	// Matrix is stored transposed,
	// but we need to use order consistent with digest.
	h := d.x[0].ByteArray()
	copy(b[:], h[:16])
	copy(b[32:], h[16:])

	h = d.x[1].ByteArray()
	copy(b[16:], h[:16])
	copy(b[48:], h[16:])
	return
}

func mulByteSliceRightx2(c00c10 *gf127.GF127x2, c01c11 *gf127.GF127x2, n int, data *byte)
