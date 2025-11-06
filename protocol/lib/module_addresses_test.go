package lib_test

import (
	"github.com/danielvindax/vd-chain/protocol/lib"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGovModuleAddress(t *testing.T) {
	require.Equal(t, "vindax10d07y265gmmuvt4z0w9aw880jnsr700jntyflm", lib.GovModuleAddress.String())
}
