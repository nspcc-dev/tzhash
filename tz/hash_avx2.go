package tz

import (
	"hash"

	"github.com/nspcc-dev/tzhash/gf127"
)

type digest2 struct {
	x [2]gf127.GF127x2
}

var _ hash.Hash = new(digest2)

func (d *digest2) Write(data []byte) (n int, err error) {
	n = len(data)
	for _, b := range data {
		mulBitRightx2(&d.x[0], &d.x[1], &minmax[(b>>7)&1])
		mulBitRightx2(&d.x[0], &d.x[1], &minmax[(b>>6)&1])
		mulBitRightx2(&d.x[0], &d.x[1], &minmax[(b>>5)&1])
		mulBitRightx2(&d.x[0], &d.x[1], &minmax[(b>>4)&1])
		mulBitRightx2(&d.x[0], &d.x[1], &minmax[(b>>3)&1])
		mulBitRightx2(&d.x[0], &d.x[1], &minmax[(b>>2)&1])
		mulBitRightx2(&d.x[0], &d.x[1], &minmax[(b>>1)&1])
		mulBitRightx2(&d.x[0], &d.x[1], &minmax[(b>>0)&1])
	}
	return
}

func (d *digest2) Sum(in []byte) []byte {
	// Make a copy of d so that caller can keep writing and summing.
	d0 := *d
	h := d0.checkSum()
	return append(in, h[:]...)
}
func (d *digest2) Reset() {
	d.x[0] = gf127.GF127x2{1, 0, 0, 0}
	d.x[1] = gf127.GF127x2{0, 0, 0, 1}
}
func (d *digest2) Size() int      { return hashSize }
func (d *digest2) BlockSize() int { return hashBlockSize }
func (d *digest2) checkSum() (b [hashSize]byte) {
	// Matrix is stored transposed,
	// but we need to use order consistent with digest.
	h := d.x[0].ByteArray()
	copy(b[:], h[:8])
	copy(b[16:], h[8:])

	h = d.x[1].ByteArray()
	copy(b[8:], h[:8])
	copy(b[24:], h[8:])
	return
}

func mulBitRightx2(c00c10 *gf127.GF127x2, c01c11 *gf127.GF127x2, e *gf127.GF127)
