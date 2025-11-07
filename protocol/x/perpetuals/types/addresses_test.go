package types_test

import (
	"testing"

	"github.com/danielvindax/vd-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestInsuranceFundModuleAddress(t *testing.T) {
	require.Equal(t, "vindax1c7ptc87hkd54e3r7zjy92q29xkq7t79w69fh2l", types.InsuranceFundModuleAddress.String())
}
