package keeper

import (
	"github.com/danielvindax/vd-chain/protocol/x/delaymsg/types"
)

var _ types.QueryServer = Keeper{}
