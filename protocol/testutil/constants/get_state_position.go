package constants

import (
	"math/big"

	clobtypes "github.com/danielvindax/vd-chain/protocol/x/clob/types"
	satypes "github.com/danielvindax/vd-chain/protocol/x/subaccounts/types"
)

var (
	// Get state position functions.
	GetStatePosition_ZeroPositionSize = func(
		subaccountId satypes.SubaccountId,
		clobPairId clobtypes.ClobPairId,
	) (
		statePositionSize *big.Int,
	) {
		return big.NewInt(0)
	}
)
