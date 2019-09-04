package tz

import (
	"github.com/nspcc-dev/tzhash/gogf127"
)

type digestp struct {
	x [4]gogf127.GF127
}

// New returns a new hash.Hash computing the Tillich-Zémor checksum.
func newPure() *digestp {
	d := new(digestp)
	d.Reset()
	return d
}

func (d *digestp) Sum(in []byte) []byte {
	// Make a copy of d so that caller can keep writing and summing.
	d0 := *d
	h := d0.checkSum()
	return append(in, h[:]...)
}

func (d *digestp) checkSum() [hashSize]byte {
	return d.byteArray()
}

func (d *digestp) byteArray() (b [hashSize]byte) {
	for i := 0; i < 4; i++ {
		t := d.x[i].ByteArray()
		copy(b[i*16:], t[:])
	}
	return
}

func (d *digestp) Reset() {
	d.x[0] = gogf127.GF127{1, 0}
	d.x[1] = gogf127.GF127{0, 0}
	d.x[2] = gogf127.GF127{0, 0}
	d.x[3] = gogf127.GF127{1, 0}
}

func (d *digestp) Write(data []byte) (n int, err error) {
	n = len(data)
	tmp := new(gogf127.GF127)
	for _, b := range data {
		mulBitRightPure(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x80 != 0, tmp)
		mulBitRightPure(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x40 != 0, tmp)
		mulBitRightPure(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x20 != 0, tmp)
		mulBitRightPure(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x10 != 0, tmp)
		mulBitRightPure(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x08 != 0, tmp)
		mulBitRightPure(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x04 != 0, tmp)
		mulBitRightPure(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x02 != 0, tmp)
		mulBitRightPure(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x01 != 0, tmp)
	}
	return
}

func (d *digestp) Size() int {
	return hashSize
}

func (d *digestp) BlockSize() int {
	return hashBlockSize
}

func mulBitRightPure(c00, c01, c10, c11 *gogf127.GF127, bit bool, tmp *gogf127.GF127) {
	if bit {
		*tmp = *c00
		gogf127.Mul10(c00, c00)
		gogf127.Add(c00, c01, c00)
		gogf127.Mul11(tmp, tmp)
		gogf127.Add(c01, tmp, c01)

		*tmp = *c10
		gogf127.Mul10(c10, c10)
		gogf127.Add(c10, c11, c10)
		gogf127.Mul11(tmp, tmp)
		gogf127.Add(c11, tmp, c11)
	} else {
		*tmp = *c00
		gogf127.Mul10(c00, c00)
		gogf127.Add(c00, c01, c00)
		*c01 = *tmp

		*tmp = *c10
		gogf127.Mul10(c10, c10)
		gogf127.Add(c10, c11, c10)
		*c11 = *tmp
	}
}
