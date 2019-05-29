package tz

import (
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	var (
		c1, c2    sl2
		n         int
		err       error
		h, h1, h2 [hashSize]byte
		b         []byte
	)

	b = make([]byte, 64)
	n, err = rand.Read(b)
	require.Equal(t, 64, n)
	require.NoError(t, err)

	// Test if our hashing is really homomorphic
	h = Sum(b)
	h1 = Sum(b[:32])
	h2 = Sum(b[32:])

	err = c1.UnmarshalBinary(h1[:])
	require.NoError(t, err)
	err = c2.UnmarshalBinary(h2[:])
	require.NoError(t, err)

	c1.Mul(&c1, &c2)
	require.Equal(t, h, c1.ByteArray())
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
		actual, expect []byte
		ps             [][]byte
		err            error
	)

	for _, tc := range testCases {
		expect, err = hex.DecodeString(tc.Hash)
		require.NoError(t, err)

		ps = make([][]byte, len(tc.Parts))
		for j := 0; j < len(tc.Parts); j++ {
			ps[j], err = hex.DecodeString(tc.Parts[j])
			require.NoError(t, err)
		}

		actual, err = Concat(ps)
		require.NoError(t, err)
		require.Equal(t, expect, actual)
	}
}

func TestValidate(t *testing.T) {
	var (
		hash []byte
		ps   [][]byte
		got  bool
		err  error
	)

	for _, tc := range testCases {
		hash, _ = hex.DecodeString(tc.Hash)
		require.NoError(t, err)

		ps = make([][]byte, len(tc.Parts))
		for j := 0; j < len(tc.Parts); j++ {
			ps[j], _ = hex.DecodeString(tc.Parts[j])
			require.NoError(t, err)
		}

		got, err = Validate(hash, ps)
		require.NoError(t, err)
		require.True(t, got)
	}
}

var testCasesSubtract = []struct {
	first, second, result string
}{
	{
		first:  "4275945919296224acd268456be23b8b2df931787a46716477e32cd991e98074029d4f03a0fedc09125ee4640d228d7d40d430659a0b2b70e9cd4d4c5361865a",
		second: "277c10e0d7c52fcc0b23ba7dbf2c3dde7dcfc1f7c0cc0d998b2de504b8c1e17c6f65ab1294aea676d4060ed2ca18c1c26fd7cec5012ab69a4ddb5e6555ac8a59",
		result: "7f5c9280352a8debea738a74abd4ec787f2c5e556800525692f651087442f9883bb97a2c1bc72d12ba26e3df8dc0f670564292ebc984976a8e353ff69a5fb3cb",
	},
	{
		first:  "18e2ce290cc74998ebd0bef76454b52a40428f13bb612e40b5b96187e9cc813248a0ed5f7ec9fb205d55d3f243e2211363f171b19eb8acc7931cf33853a79069",
		second: "73a0582fa7d00d62fd09c1cd18589cdb2b126cb58b3a022ae47a8a787dabe35c4388aaf0d8bb343b1e58ee8d267812d115f40a0da611f42458f452e102f60700",
		result: "54ccaad1bb15b2989fa31109713bca955ea5d87bbd3113b3008cea167c00052266e9c9fcb73ece98c6c08cccb074ba3d39b5d8685f022fc388e2bf1997c5bd1d",
	},
}

func TestSubtract(t *testing.T) {
	var (
		a, b, c, r []byte
		err        error
	)

	for _, tc := range testCasesSubtract {
		a, err = hex.DecodeString(tc.first)
		require.NoError(t, err)

		b, err = hex.DecodeString(tc.second)
		require.NoError(t, err)

		c, err = hex.DecodeString(tc.result)
		require.NoError(t, err)

		r, err = SubtractR(c, b)
		require.NoError(t, err)
		require.Equal(t, a, r)

		r, err = SubtractL(c, a)
		require.NoError(t, err)
		require.Equal(t, b, r)
	}
}
