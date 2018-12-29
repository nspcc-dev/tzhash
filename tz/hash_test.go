package tz

import (
	"encoding/hex"
	"math/rand"
	"testing"

	. "github.com/onsi/gomega"
)

func TestHash(t *testing.T) {
	var (
		c1, c2    sl2
		n         int
		err       error
		h, h1, h2 [hashSize]byte
		b         []byte
	)

	g := NewGomegaWithT(t)

	b = make([]byte, 64)
	n, err = rand.Read(b)
	g.Expect(n).To(Equal(64))
	g.Expect(err).NotTo(HaveOccurred())

	// Test if our hashing is really homomorphic
	h = Sum(b)
	h1 = Sum(b[:32])
	h2 = Sum(b[32:])

	err = c1.UnmarshalBinary(h1[:])
	g.Expect(err).NotTo(HaveOccurred())
	err = c2.UnmarshalBinary(h2[:])
	g.Expect(err).NotTo(HaveOccurred())

	c1.Mul(&c1, &c2)
	g.Expect(c1.ByteArray()).To(Equal(h))
}

var testCases = []struct {
	Hash  string
	Parts []string
}{{
	Hash: "7f5c9280352a8debea738a74abd4ec787f2c5e556800525692f651087442f9883bb97a2c1bc72d12ba26e3df8dc0f670564292ebc984976a8e353ff69a5fb3cb",
	Parts: []string{
		"4275945919296224acd268456be23b8b2df931787a46716477e32cd991e98074029d4f03a0fedc09125ee4640d228d7d40d430659a0b2b70e9cd4d4c5361865a",
		"2828661d1b1e77f21788d3b365f140a2395d57dc2083c33e60d9a80e69017d5016a249c7adfe1718a10ba887dedbdaec5c4c1fbecdb1f98776b43f1142c26a88",
		"02310598b45dfa77db9f00eed6ab60773dd8bed7bdac431b42e441fae463f64c6e2688402cfdcec5def47a299b0651fb20878cf4410991bd57056d7b4b31635a",
		"1ed7e0b065c060d915e7355cdcb4edc752c06d2a4b39d90c8985aeb58e08cb9e5bbe4b2b45524efbd68cd7e4081a1b8362941200a4c9f76a0a9f9ac9b7868c03",
		"6f11e3dc4fff99ffa45e36e4655cfc657c29e950e598a90f426bf5710de9171323523db7636643b23892783f4fb3cf8e583d584c82d29558a105a615a668fc9e",
		"1865dbdb4c849620fb2c4809d75d62490f83c11f2145abaabbdc9a66ae58ce1f2e42c34d3b380e5dea1b45217750b42d130f995b162afbd2e412b0d41ec8871b",
		"5102dd1bd1f08f44dbf3f27ac895020d63f96044ce3b491aed3efbc7bbe363bc5d800101d63890f89a532427812c30c9674f37476ba44daf758afa88d4f91063",
		"70cab735dad90164cc61f7411396221c4e549f12392c0d77728c89a9754f606c7d961169d4fa88133a1ba954bad616656c86f8fd1335a2f3428fd4dca3a3f5a5",
		"430f3e92536ff9a50cbcdf08d8810a59786ca37e31d54293646117a93469f61c6cdd67933128407d77f3235293293ee86dbc759d12dfe470969eba1b4a373bd0",
		"46e1d97912ca2cf92e6a9a63667676835d900cdb2fff062136a64d8d60a8e5aa644ccee3558900af8e77d56b013ed5da12d9d0b7de0f56976e040b3d01345c0d",
	},
}}

func TestConcat(t *testing.T) {
	var (
		got, expected []byte
		ps            [][]byte
		err           error
	)

	g := NewGomegaWithT(t)

	for _, tc := range testCases {
		expected, err = hex.DecodeString(tc.Hash)
		g.Expect(err).NotTo(HaveOccurred())

		ps = make([][]byte, len(tc.Parts))
		for j := 0; j < len(tc.Parts); j++ {
			ps[j], err = hex.DecodeString(tc.Parts[j])
			g.Expect(err).NotTo(HaveOccurred())
		}

		got, err = Concat(ps)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(got).To(Equal(expected))
	}
}

func TestValidate(t *testing.T) {
	var (
		hash []byte
		ps   [][]byte
		got  bool
		err  error
	)

	g := NewGomegaWithT(t)

	for _, tc := range testCases {
		hash, _ = hex.DecodeString(tc.Hash)
		g.Expect(err).NotTo(HaveOccurred())

		ps = make([][]byte, len(tc.Parts))
		for j := 0; j < len(tc.Parts); j++ {
			ps[j], _ = hex.DecodeString(tc.Parts[j])
			g.Expect(err).NotTo(HaveOccurred())
		}

		got, err = Validate(hash, ps)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(got).To(Equal(true))
	}
}
