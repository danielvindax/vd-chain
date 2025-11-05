package keeper

import (
	"github.com/danielvindax/vd-chain/protocol/x/revshare/types"
)

var _ types.QueryServer = Keeper{}
