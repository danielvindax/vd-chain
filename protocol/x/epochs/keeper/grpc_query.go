package keeper

import (
	"github.com/danielvindax/vd-chain/protocol/x/epochs/types"
)

var _ types.QueryServer = Keeper{}
