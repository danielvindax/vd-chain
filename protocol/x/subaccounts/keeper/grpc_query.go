package keeper

import (
	"github.com/danielvindax/vd-chain/protocol/x/subaccounts/types"
)

var _ types.QueryServer = Keeper{}
