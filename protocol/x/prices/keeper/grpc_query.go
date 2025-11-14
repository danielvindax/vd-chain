package keeper

import (
	"github.com/danielvindax/vd-chain/protocol/x/prices/types"
)

var _ types.QueryServer = Keeper{}
