package types_test

import (
	"github.com/danielvindax/vd-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTreasuryModuleAddress(t *testing.T) {
	require.Equal(t, "vindax16wrau2x4tsg033xfrrdpae6kxfn9kyuernd6m7", types.TreasuryModuleAddress.String())
}
