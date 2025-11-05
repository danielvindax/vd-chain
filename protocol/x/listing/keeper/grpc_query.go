package keeper

import (
	"github.com/danielvindax/vd-chain/protocol/x/listing/types"
)

var _ types.QueryServer = Keeper{}
