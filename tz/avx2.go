// Copyright 2019 (c) NSPCC
//
// This file contains AVX2 implementation.
package tz

import (
	"hash"

	"github.com/nspcc-dev/tzhash/gf127"
	"github.com/nspcc-dev/tzhash/gf127/avx2"
)

type digest2 struct {
	x [2]avx2.GF127x2
}

// type assertion
var _ hash.Hash = (*digest2)(nil)

func newAVX2() *digest2 {
	d := new(digest2)
	d.Reset()
	return d
}

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
	d.x[0] = avx2.GF127x2{gf127.GF127{1, 0}, gf127.GF127{0, 0}}
	d.x[1] = avx2.GF127x2{gf127.GF127{0, 0}, gf127.GF127{1, 0}}
}
func (d *digest2) Size() int      { return hashSize }
func (d *digest2) BlockSize() int { return hashBlockSize }
func (d *digest2) checkSum() (b [hashSize]byte) {
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

func mulBitRightx2(c00c10 *avx2.GF127x2, c01c11 *avx2.GF127x2, e *gf127.GF127)
