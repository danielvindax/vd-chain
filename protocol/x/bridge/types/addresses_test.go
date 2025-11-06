package types_test

import (
	"github.com/danielvindax/vd-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModuleAddress(t *testing.T) {
	require.Equal(t, "vindax1zlefkpe3g0vvm9a4h0jf9000lmqutlh9j7tmen", types.ModuleAddress.String())
}
