package tz

import (
	"errors"

	"github.com/nspcc-dev/tzhash/gf127"
	"github.com/nspcc-dev/tzhash/gf127/avx"
)

type sl2 [2][2]gf127.GF127

var id = sl2{
	{gf127.GF127{1, 0}, gf127.GF127{0, 0}},
	{gf127.GF127{0, 0}, gf127.GF127{1, 0}},
}

var mul func(a, b, c *sl2, x *[4]gf127.GF127)

func init() {
	if hasAVX {
		mul = mulSL2AVX
	} else {
		mul = mulSL2Pure
	}
}

func (c *sl2) MarshalBinary() (data []byte, err error) {
	s := c.ByteArray()
	return s[:], nil
}

func (c *sl2) UnmarshalBinary(data []byte) (err error) {
	if len(data) != 64 {
		return errors.New("data must be 64-bytes long")
	}

	if err = c[0][0].UnmarshalBinary(data[:16]); err != nil {
		return
	}
	if err = c[0][1].UnmarshalBinary(data[16:32]); err != nil {
		return
	}
	if err = c[1][0].UnmarshalBinary(data[32:48]); err != nil {
		return
	}
	if err = c[1][1].UnmarshalBinary(data[48:64]); err != nil {
		return
	}

	return
}

func (c *sl2) mulStrassen(a, b *sl2, x *[8]gf127.GF127) *sl2 {
	// strassen algorithm
	avx.Add(&a[0][0], &a[1][1], &x[0])
	avx.Add(&b[0][0], &b[1][1], &x[1])
	avx.Mul(&x[0], &x[1], &x[0])

	avx.Add(&a[1][0], &a[1][1], &x[1])
	avx.Mul(&x[1], &b[0][0], &x[1])

	avx.Add(&b[0][1], &b[1][1], &x[2])
	avx.Mul(&x[2], &a[0][0], &x[2])

	avx.Add(&b[1][0], &b[0][0], &x[3])
	avx.Mul(&x[3], &a[1][1], &x[3])

	avx.Add(&a[0][0], &a[0][1], &x[4])
	avx.Mul(&x[4], &b[1][1], &x[4])

	avx.Add(&a[1][0], &a[0][0], &x[5])
	avx.Add(&b[0][0], &b[0][1], &x[6])
	avx.Mul(&x[5], &x[6], &x[5])

	avx.Add(&a[0][1], &a[1][1], &x[6])
	avx.Add(&b[1][0], &b[1][1], &x[7])
	avx.Mul(&x[6], &x[7], &x[6])

	avx.Add(&x[2], &x[4], &c[0][1])
	avx.Add(&x[1], &x[3], &c[1][0])

	avx.Add(&x[4], &x[6], &x[4])
	avx.Add(&x[0], &x[3], &c[0][0])
	avx.Add(&c[0][0], &x[4], &c[0][0])

	avx.Add(&x[0], &x[1], &x[0])
	avx.Add(&x[2], &x[5], &c[1][1])
	avx.Add(&c[1][1], &x[0], &c[1][1])

	return c
}

func mulSL2AVX(a, b, c *sl2, x *[4]gf127.GF127) {
	avx.Mul(&a[0][0], &b[0][0], &x[0])
	avx.Mul(&a[0][0], &b[0][1], &x[1])
	avx.Mul(&a[1][0], &b[0][0], &x[2])
	avx.Mul(&a[1][0], &b[0][1], &x[3])

	avx.Mul(&a[0][1], &b[1][0], &c[0][0])
	avx.Add(&c[0][0], &x[0], &c[0][0])
	avx.Mul(&a[0][1], &b[1][1], &c[0][1])
	avx.Add(&c[0][1], &x[1], &c[0][1])
	avx.Mul(&a[1][1], &b[1][0], &c[1][0])
	avx.Add(&c[1][0], &x[2], &c[1][0])
	avx.Mul(&a[1][1], &b[1][1], &c[1][1])
	avx.Add(&c[1][1], &x[3], &c[1][1])
}

func mulSL2Pure(a, b, c *sl2, x *[4]gf127.GF127) {
	gf127.Mul((*gf127.GF127)(&a[0][0]), (*gf127.GF127)(&b[0][0]), (*gf127.GF127)(&x[0]))
	gf127.Mul((*gf127.GF127)(&a[0][0]), (*gf127.GF127)(&b[0][1]), (*gf127.GF127)(&x[1]))
	gf127.Mul((*gf127.GF127)(&a[1][0]), (*gf127.GF127)(&b[0][0]), (*gf127.GF127)(&x[2]))
	gf127.Mul((*gf127.GF127)(&a[1][0]), (*gf127.GF127)(&b[0][1]), (*gf127.GF127)(&x[3]))

	gf127.Mul((*gf127.GF127)(&a[0][1]), (*gf127.GF127)(&b[1][0]), (*gf127.GF127)(&c[0][0]))
	gf127.Add((*gf127.GF127)(&c[0][0]), (*gf127.GF127)(&x[0]), (*gf127.GF127)(&c[0][0]))
	gf127.Mul((*gf127.GF127)(&a[0][1]), (*gf127.GF127)(&b[1][1]), (*gf127.GF127)(&c[0][1]))
	gf127.Add((*gf127.GF127)(&c[0][1]), (*gf127.GF127)(&x[1]), (*gf127.GF127)(&c[0][1]))
	gf127.Mul((*gf127.GF127)(&a[1][1]), (*gf127.GF127)(&b[1][0]), (*gf127.GF127)(&c[1][0]))
	gf127.Add((*gf127.GF127)(&c[1][0]), (*gf127.GF127)(&x[2]), (*gf127.GF127)(&c[1][0]))
	gf127.Mul((*gf127.GF127)(&a[1][1]), (*gf127.GF127)(&b[1][1]), (*gf127.GF127)(&c[1][1]))
	gf127.Add((*gf127.GF127)(&c[1][1]), (*gf127.GF127)(&x[3]), (*gf127.GF127)(&c[1][1]))
}

func (c *sl2) MulA() *sl2 {
	var a gf127.GF127

	avx.Mul10(&c[0][0], &a)
	avx.Mul1(&c[0][0], &c[0][1])
	avx.Add(&a, &c[0][1], &c[0][0])

	avx.Mul10(&c[1][0], &a)
	avx.Mul1(&c[1][0], &c[1][1])
	avx.Add(&a, &c[1][1], &c[1][0])

	return c
}

func (c *sl2) MulB() *sl2 {
	var a gf127.GF127

	avx.Mul1(&c[0][0], &a)
	avx.Mul10(&c[0][0], &c[0][0])
	avx.Add(&c[0][1], &c[0][0], &c[0][0])
	avx.Add(&c[0][0], &a, &c[0][1])

	avx.Mul1(&c[1][0], &a)
	avx.Mul10(&c[1][0], &c[1][0])
	avx.Add(&c[1][1], &c[1][0], &c[1][0])
	avx.Add(&c[1][0], &a, &c[1][1])

	return c
}

func (c *sl2) Mul(a, b *sl2) *sl2 {
	mul(a, b, c, new([4]gf127.GF127))
	return c
}

// Inv returns inverse of a in GL_2(GF(2^127))
func Inv(a *sl2) (b *sl2) {
	b = new(sl2)
	inv(a, b, new([2]gf127.GF127))
	return
}

func inv(a, b *sl2, t *[2]gf127.GF127) {
	avx.Mul(&a[0][0], &a[1][1], &t[0])
	avx.Mul(&a[0][1], &a[1][0], &t[1])
	avx.Add(&t[0], &t[1], &t[0])
	avx.Inv(&t[0], &t[1])

	avx.Mul(&t[1], &a[0][0], &b[1][1])
	avx.Mul(&t[1], &a[0][1], &b[0][1])
	avx.Mul(&t[1], &a[1][0], &b[1][0])
	avx.Mul(&t[1], &a[1][1], &b[0][0])
}

func (c *sl2) String() string {
	return c[0][0].String() + c[0][1].String() +
		c[1][0].String() + c[1][1].String()
}

func (c *sl2) ByteArray() (b [hashSize]byte) {
	t := c[0][0].ByteArray()
	copy(b[:], t)

	t = c[0][1].ByteArray()
	copy(b[16:], t)

	t = c[1][0].ByteArray()
	copy(b[32:], t)

	t = c[1][1].ByteArray()
	copy(b[48:], t)

	return
}
