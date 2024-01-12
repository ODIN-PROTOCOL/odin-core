package odinrng_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ODIN-PROTOCOL/odin-core/pkg/odinrng"
)

func TestRngRandom(t *testing.T) {
	r, err := odinrng.NewRng([]byte("THIS_IS_A_RANDOM_SEED_LONG_ENOUGH_FOR_ENTROPY"), []byte("1"), []byte("TEST"))
	require.NoError(t, err)
	require.Equal(t, uint64(5751621621077249396), r.NextUint64())
	require.Equal(t, uint64(16474548556352052882), r.NextUint64())
	require.Equal(t, uint64(17097048712898369316), r.NextUint64())
	require.Equal(t, uint64(10686498023352306525), r.NextUint64())
	require.Equal(t, uint64(2144097648649487685), r.NextUint64())
	require.Equal(t, uint64(1642256529570429276), r.NextUint64())
	require.Equal(t, uint64(1298883664373060799), r.NextUint64())
}

func TestRngRandomNotEnoughEntropy(t *testing.T) {
	_, err := odinrng.NewRng([]byte("TOO_SHORT"), []byte("1"), []byte("TEST"))
	require.EqualError(t, err, "drbg: insufficient entropyInput")
}
