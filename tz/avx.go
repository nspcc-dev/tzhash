// Copyright 2018 (c) NSPCC
//
// This file contains AVX implementation.
package tz

import (
	"hash"
	"math"
)

type digest struct {
	x [4]GF127
}

// type assertion
var _ hash.Hash = (*digest)(nil)

var (
	minmax  = [2]GF127{{0, 0}, {math.MaxUint64, math.MaxUint64}}
	x127x63 = GF127{1 << 63, 1 << 63} //nolint:deadcode,varcheck
)

func newAVX() *digest {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Sum(in []byte) []byte {
	// Make a copy of d so that caller can keep writing and summing.
	d0 := *d
	h := d0.checkSum()
	return append(in, h[:]...)
}

func (d *digest) checkSum() [hashSize]byte {
	return d.byteArray()
}

func (d *digest) byteArray() (b [hashSize]byte) {
	copy(b[:], d.x[0].ByteArray())
	copy(b[16:], d.x[1].ByteArray())
	copy(b[32:], d.x[2].ByteArray())
	copy(b[48:], d.x[3].ByteArray())
	return
}

func (d *digest) Reset() {
	d.x[0] = GF127{1, 0}
	d.x[1] = GF127{0, 0}
	d.x[2] = GF127{0, 0}
	d.x[3] = GF127{1, 0}
}

func (d *digest) Write(data []byte) (n int, err error) {
	n = len(data)
	for _, b := range data {
		mulBitRight(&d.x[0], &d.x[1], &d.x[2], &d.x[3], &minmax[(b>>7)&1])
		mulBitRight(&d.x[0], &d.x[1], &d.x[2], &d.x[3], &minmax[(b>>6)&1])
		mulBitRight(&d.x[0], &d.x[1], &d.x[2], &d.x[3], &minmax[(b>>5)&1])
		mulBitRight(&d.x[0], &d.x[1], &d.x[2], &d.x[3], &minmax[(b>>4)&1])
		mulBitRight(&d.x[0], &d.x[1], &d.x[2], &d.x[3], &minmax[(b>>3)&1])
		mulBitRight(&d.x[0], &d.x[1], &d.x[2], &d.x[3], &minmax[(b>>2)&1])
		mulBitRight(&d.x[0], &d.x[1], &d.x[2], &d.x[3], &minmax[(b>>1)&1])
		mulBitRight(&d.x[0], &d.x[1], &d.x[2], &d.x[3], &minmax[(b>>0)&1])
	}
	return
}

func (d *digest) Size() int {
	return hashSize
}

func (d *digest) BlockSize() int {
	return hashBlockSize
}

func mulBitRight(c00, c01, c10, c11, e *GF127)
