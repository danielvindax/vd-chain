package keeper

import (
	"github.com/danielvindax/vd-chain/protocol/x/perpetuals/types"
)

var _ types.QueryServer = Keeper{}
