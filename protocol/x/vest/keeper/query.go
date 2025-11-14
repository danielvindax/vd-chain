package keeper

import (
	"github.com/danielvindax/vd-chain/protocol/x/vest/types"
)

var _ types.QueryServer = Keeper{}
