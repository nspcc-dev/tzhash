package tz

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func newDigest() *digest {
	d := &digest{}
	d.Reset()
	return d
}

func Test_digest_marshaling(t *testing.T) {
	var (
		d              = newDigest()
		marshaledState []byte
		hashSum        []byte
	)

	for i := byte(0); i < 10; i++ {
		n, err := d.Write([]byte{i})
		require.NoError(t, err)
		require.Equal(t, 1, n)

		t.Run("marshal", func(t *testing.T) {
			marshaledState, err = d.MarshalBinary()

			require.NoError(t, err)
			require.Len(t, marshaledState, Size)

			hashSum = d.Sum(nil)
			require.Len(t, hashSum, Size)
		})

		t.Run("unmarshal", func(t *testing.T) {
			unmarshalDigest := newDigest()
			err = unmarshalDigest.UnmarshalBinary(marshaledState)
			require.NoError(t, err)

			unmarshalDigestHash := unmarshalDigest.Sum(nil)
			require.Len(t, unmarshalDigestHash, Size)

			require.Equal(t, hashSum, unmarshalDigestHash)
		})
	}

	t.Run("invalid length", func(t *testing.T) {
		unmarshalDigest := newDigest()
		state := []byte{1, 2, 3}

		err := unmarshalDigest.UnmarshalBinary(state)
		require.Error(t, err)
	})

	t.Run("invalid state data", func(t *testing.T) {
		someDigest := newDigest()
		_, err := d.Write([]byte{1, 2, 3})
		require.NoError(t, err)

		state, err := someDigest.MarshalBinary()
		require.NoError(t, err)

		a := uint64(1) << 63
		// broke state data.
		binary.BigEndian.PutUint64(state[0:8], a)
		unmarshalDigest := newDigest()

		err = unmarshalDigest.UnmarshalBinary(state)
		require.Error(t, err)
	})
}
