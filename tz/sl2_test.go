package tz

import (
	"math/rand"
	"testing"
	"time"

	"github.com/nspcc-dev/tzhash/gf127"
	. "github.com/onsi/gomega"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func random() (a *sl2) {
	a = new(sl2)
	a[0][0] = *gf127.Random()
	a[0][1] = *gf127.Random()
	a[1][0] = *gf127.Random()

	// so that result is in SL2
	// d = a^-1*(1+b*c)
	gf127.Mul(&a[0][1], &a[1][0], &a[1][1])
	gf127.Add(&a[1][1], gf127.New(1, 0), &a[1][1])

	t := gf127.New(0, 0)
	gf127.Inv(&a[0][0], t)
	gf127.Mul(t, &a[1][1], &a[1][1])

	return
}

func TestSL2_MarshalBinary(t *testing.T) {
	var (
		a = random()
		b = new(sl2)
		g = NewGomegaWithT(t)
	)

	data, err := a.MarshalBinary()
	g.Expect(err).NotTo(HaveOccurred())

	err = b.UnmarshalBinary(data)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(a).To(Equal(b))
}

func TestInv(t *testing.T) {
	var (
		a, b, c *sl2
		g       = NewGomegaWithT(t)
	)

	c = new(sl2)
	for i := 0; i < 5; i++ {
		a = random()
		b = Inv(a)
		c = c.Mul(a, b)

		g.Expect(*c).To(Equal(id))
	}
}
