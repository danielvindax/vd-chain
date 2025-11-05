package keeper

import (
	"github.com/danielvindax/vd-chain/protocol/x/assets/types"
)

var _ types.QueryServer = Keeper{}
