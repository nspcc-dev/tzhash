// Copyright 2018 (c) NSPCC
//
// This file contains AVX implementation.
package tz

import (
	"hash"
)

type digest4 struct {
	x [4]GF127
}

// type assertion
var _ hash.Hash = (*digest4)(nil)

func newAVXInline() *digest4 {
	d := new(digest4)
	d.Reset()
	return d
}

func (d *digest4) Sum(in []byte) []byte {
	// Make a copy of d so that caller can keep writing and summing.
	d0 := *d
	h := d0.checkSum()
	return append(in, h[:]...)
}

func (d *digest4) checkSum() [Size]byte {
	return d.byteArray()
}

func (d *digest4) byteArray() (b [Size]byte) {
	copy(b[:], d.x[0].ByteArray())
	copy(b[16:], d.x[1].ByteArray())
	copy(b[32:], d.x[2].ByteArray())
	copy(b[48:], d.x[3].ByteArray())
	return
}

func (d *digest4) Reset() {
	d.x[0] = GF127{1, 0}
	d.x[1] = GF127{0, 0}
	d.x[2] = GF127{0, 0}
	d.x[3] = GF127{1, 0}
}

func (d *digest4) Write(data []byte) (n int, err error) {
	n = len(data)
	for _, b := range data {
		mulByteRight(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b)
	}
	return
}

func (d *digest4) Size() int {
	return Size
}

func (d *digest4) BlockSize() int {
	return hashBlockSize
}

func mulByteRight(c00, c01, c10, c11 *GF127, b byte)
