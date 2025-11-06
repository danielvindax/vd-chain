package types_test

import (
	"github.com/danielvindax/vd-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModuleAddress(t *testing.T) {
	require.Equal(t, "vindax1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5vcx39", types.ModuleAddress.String())
}
